package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

// Reader for random number generation mocking
type mockReader struct {
	err bool
}

func (m *mockReader) Read(p []byte) (n int, err error) {
	if m.err {
		return 0, errors.New("mock rand error")
	}
	for i := range p {
		p[i] = byte(i)
	}
	return len(p), nil
}

// Mocks cipher functions from crypto/cipher
type MockCipherFactory struct {
	cipherError bool
	gcmError    bool
}

func (m *MockCipherFactory) newCipher(key []byte) (cipher.Block, error) {
	if m.cipherError {
		return nil, fmt.Errorf("some error creating Cipher")
	}
	return aes.NewCipher(make([]byte, 16))
}

func (m *MockCipherFactory) newGCM(block cipher.Block) (cipher.AEAD, error) {
	if m.gcmError {
		return nil, fmt.Errorf("some error creating GCM")
	}
	return cipher.NewGCM(block)
}

func TestGenerateAES(t *testing.T) {
	type args struct {
		size       int
		randReader Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *AES
		wantErr bool
	}{
		{
			name: "Valid size generates key",
			args: args{
				size:       16,
				randReader: &mockReader{},
			},
			wantErr: false,
		},
		{
			name: "Invalid size returns error",
			args: args{
				size:       8,
				randReader: &mockReader{},
			},
			wantErr: true,
		},
		{
			name: "Returns error on failing to generate random key",
			args: args{
				size:       16,
				randReader: &mockReader{err: true},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateAES(tt.args.size, tt.args.randReader)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateAES() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && len(got.GetKey()) != tt.args.size {
				t.Errorf("GenerateAES() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncryptWithAESGCM(t *testing.T) {

	type args struct {
		factory   CipherFactory
		plaintext []byte
	}
	type fields struct {
		size   int
		reader Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Encrypts message successfully",
			fields: fields{
				size:   16,
				reader: &mockReader{err: false},
			},
			args: args{
				factory:   &MockCipherFactory{},
				plaintext: []byte("test_message"),
			},
			wantErr: false,
		},
		{
			name: "Returns error on failing to generate nonce",
			fields: fields{
				size:   16,
				reader: &mockReader{err: true},
			},
			args: args{
				factory:   &MockCipherFactory{},
				plaintext: []byte("test_message"),
			},
			wantErr: true,
		},
		{
			name: "Returns error on failing to create new cipher",
			fields: fields{
				size:   16,
				reader: &mockReader{},
			},
			args: args{
				factory:   &MockCipherFactory{cipherError: true},
				plaintext: []byte("test_message"),
			},
			wantErr: true,
		},
		{
			name: "Returns error on failing to create new GCM",
			fields: fields{
				size:   16,
				reader: &mockReader{},
			},
			args: args{
				factory:   &MockCipherFactory{gcmError: true},
				plaintext: []byte("test_message"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, _ := GenerateAES(tt.fields.size, tt.fields.reader)
			_, err := a.EncryptWithAESGCM(tt.args.factory, tt.fields.reader, tt.args.plaintext)
			if (err != nil) != tt.wantErr {
				t.Errorf("AES.EncryptWithAESGCM() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			rand.Reader = &mockReader{}
		})
	}
}

func TestAES_DecryptWithAESGCM(t *testing.T) {
	randReader := &mockReader{}
	aesTest, _ := GenerateAES(16, randReader)

	ciphertext, _ := aesTest.EncryptWithAESGCM(&MockCipherFactory{}, randReader, []byte("test_message"))

	type args struct {
		factory    CipherFactory
		ciphertext []byte
	}
	tests := []struct {
		name          string
		args          args
		wantPlaintext []byte
		wantErr       bool
	}{
		{
			name: "Decrypts message successfully",
			args: args{
				factory:    &MockCipherFactory{},
				ciphertext: ciphertext,
			},
			wantPlaintext: []byte("test_message"),
			wantErr:       false,
		},
		{
			name: "Returns error on failing to create new Cipher",
			args: args{
				factory:    &MockCipherFactory{cipherError: true},
				ciphertext: ciphertext,
			},
			wantPlaintext: nil,
			wantErr:       true,
		},
		{
			name: "Returns error on failing to create new GCM",
			args: args{
				factory:    &MockCipherFactory{gcmError: true},
				ciphertext: ciphertext,
			},
			wantPlaintext: nil,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotPlaintext, err := aesTest.DecryptWithAESGCM(tt.args.factory, tt.args.ciphertext)
			if (err != nil) != tt.wantErr {
				t.Errorf("AES.DecryptWithAESGCM() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPlaintext, tt.wantPlaintext) {
				t.Errorf("AES.DecryptWithAESGCM() = %v, want %v", gotPlaintext, tt.wantPlaintext)
			}
		})
	}
}
