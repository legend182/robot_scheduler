package utils

import (
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"errors"
)

// DESEncrypt DES加密
func DESEncrypt(plaintext, key string) (string, error) {
	// 确保密钥长度为8字节
	keyBytes := []byte(key)
	if len(keyBytes) != 8 {
		return "", errors.New("DES key must be 8 bytes")
	}

	// 创建DES cipher
	block, err := des.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// PKCS5Padding填充
	plaintextBytes := []byte(plaintext)
	plaintextBytes = pkcs5Padding(plaintextBytes, block.BlockSize())

	// 使用CBC模式
	iv := keyBytes // 使用密钥作为IV（实际生产环境应使用随机IV）
	mode := cipher.NewCBCEncrypter(block, iv)

	// 加密
	ciphertext := make([]byte, len(plaintextBytes))
	mode.CryptBlocks(ciphertext, plaintextBytes)

	// Base64编码
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DESDecrypt DES解密
func DESDecrypt(ciphertext, key string) (string, error) {
	// 确保密钥长度为8字节
	keyBytes := []byte(key)
	if len(keyBytes) != 8 {
		return "", errors.New("DES key must be 8 bytes")
	}

	// Base64解码
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// 创建DES cipher
	block, err := des.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// 使用CBC模式
	iv := keyBytes // 使用密钥作为IV（实际生产环境应使用随机IV）
	mode := cipher.NewCBCDecrypter(block, iv)

	// 解密
	plaintext := make([]byte, len(ciphertextBytes))
	mode.CryptBlocks(plaintext, ciphertextBytes)

	// 去除PKCS5Padding填充
	plaintext = pkcs5UnPadding(plaintext)

	return string(plaintext), nil
}

// pkcs5Padding PKCS5填充
func pkcs5Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	return append(data, padtext...)
}

// pkcs5UnPadding 去除PKCS5填充
func pkcs5UnPadding(data []byte) []byte {
	length := len(data)
	if length == 0 {
		return data
	}
	unpadding := int(data[length-1])
	if unpadding > length {
		return data
	}
	return data[:(length - unpadding)]
}
