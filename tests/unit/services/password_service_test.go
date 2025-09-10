package services_test

import (
	"testing"
	
	"github.com/JoseLuis21/mv-backend/internal/core/auth/services"
)

func TestPasswordHasher(t *testing.T) {
	hasher := services.NewPasswordHasher()
	
	t.Run("Hash valid password", func(t *testing.T) {
		password := "testpassword123"
		hash, err := hasher.Hash(password)
		
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		
		if hash == "" {
			t.Error("Expected hash to be generated")
		}
		
		if hash == password {
			t.Error("Hash should not be the same as original password")
		}
		
		if len(hash) < 20 {
			t.Error("Hash should be at least 20 characters long")
		}
	})
	
	t.Run("Hash empty password", func(t *testing.T) {
		_, err := hasher.Hash("")
		
		if err == nil {
			t.Error("Expected error for empty password")
		}
		
		expectedError := "contraseña no puede estar vacía"
		if err.Error() != expectedError {
			t.Errorf("Expected error %q, got %q", expectedError, err.Error())
		}
	})
	
	t.Run("Hash short password", func(t *testing.T) {
		_, err := hasher.Hash("1234567") // 7 characters
		
		if err == nil {
			t.Error("Expected error for short password")
		}
		
		expectedError := "contraseña debe tener al menos 8 caracteres"
		if err.Error() != expectedError {
			t.Errorf("Expected error %q, got %q", expectedError, err.Error())
		}
	})
	
	t.Run("Hash very long password", func(t *testing.T) {
		longPassword := make([]byte, 73) // 73 characters (bcrypt limit is 72)
		for i := range longPassword {
			longPassword[i] = 'a'
		}
		
		_, err := hasher.Hash(string(longPassword))
		
		if err == nil {
			t.Error("Expected error for very long password")
		}
		
		expectedError := "contraseña muy larga (máximo 72 caracteres)"
		if err.Error() != expectedError {
			t.Errorf("Expected error %q, got %q", expectedError, err.Error())
		}
	})
	
	t.Run("Verify correct password", func(t *testing.T) {
		password := "testpassword123"
		hash, err := hasher.Hash(password)
		
		if err != nil {
			t.Fatalf("Failed to hash password: %v", err)
		}
		
		err = hasher.Verify(hash, password)
		if err != nil {
			t.Errorf("Expected no error verifying correct password, got %v", err)
		}
	})
	
	t.Run("Verify incorrect password", func(t *testing.T) {
		password := "testpassword123"
		wrongPassword := "wrongpassword"
		
		hash, err := hasher.Hash(password)
		if err != nil {
			t.Fatalf("Failed to hash password: %v", err)
		}
		
		err = hasher.Verify(hash, wrongPassword)
		if err == nil {
			t.Error("Expected error verifying incorrect password")
		}
		
		expectedError := "contraseña incorrecta"
		if err.Error() != expectedError {
			t.Errorf("Expected error %q, got %q", expectedError, err.Error())
		}
	})
	
	t.Run("Verify with empty hash", func(t *testing.T) {
		err := hasher.Verify("", "password")
		
		if err == nil {
			t.Error("Expected error for empty hash")
		}
		
		expectedError := "hash y contraseña son requeridos"
		if err.Error() != expectedError {
			t.Errorf("Expected error %q, got %q", expectedError, err.Error())
		}
	})
	
	t.Run("Verify with empty password", func(t *testing.T) {
		hash := "$2a$12$somehash"
		err := hasher.Verify(hash, "")
		
		if err == nil {
			t.Error("Expected error for empty password")
		}
		
		expectedError := "hash y contraseña son requeridos"
		if err.Error() != expectedError {
			t.Errorf("Expected error %q, got %q", expectedError, err.Error())
		}
	})
	
	t.Run("Different hashes for same password", func(t *testing.T) {
		password := "testpassword123"
		
		hash1, err1 := hasher.Hash(password)
		hash2, err2 := hasher.Hash(password)
		
		if err1 != nil || err2 != nil {
			t.Fatalf("Failed to hash password: %v, %v", err1, err2)
		}
		
		if hash1 == hash2 {
			t.Error("Different hash calls should produce different hashes (due to salt)")
		}
		
		// But both should verify correctly
		if err := hasher.Verify(hash1, password); err != nil {
			t.Errorf("Hash1 should verify correctly: %v", err)
		}
		
		if err := hasher.Verify(hash2, password); err != nil {
			t.Errorf("Hash2 should verify correctly: %v", err)
		}
	})
}

// Benchmark para medir performance del hashing
func BenchmarkPasswordHash(b *testing.B) {
	hasher := services.NewPasswordHasher()
	password := "testpassword123"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hasher.Hash(password)
	}
}

func BenchmarkPasswordVerify(b *testing.B) {
	hasher := services.NewPasswordHasher()
	password := "testpassword123"
	hash, _ := hasher.Hash(password)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hasher.Verify(hash, password)
	}
}