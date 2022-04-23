package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"syscall/js"
)

func main() {

	c := make(chan struct{}, 0)
	js.Global().Set("encrypt", js.FuncOf(encrypt))
	js.Global().Set("decrypt", js.FuncOf(decrypt))
	js.Global().Set("getkey", js.FuncOf(getKeyPrompt))

	fmt.Println("wasm başlatıldı")
	<-c

}

func getKeyPrompt(this js.Value, inputs []js.Value) interface{} {
	data := js.Global().Get("prompt").Invoke("Please enter public key:", "")
	if len(data.String()) == 16 {
		return data
	} else {
		js.Global().Get("alert").Invoke("wrong format")
		js.Global().Call("getkey")
	}
	return "0000000000000000"
}

// AES kullanarak dizeyi base64'e çeviriyoruz
func encrypt(this js.Value, inputs []js.Value) interface{} {
	keyTobyt := []byte(inputs[0].String())
	plaintext := []byte(inputs[1].String())

	// alınan public key ile aes bloğu oluşturuyoruz
	block, err := aes.NewCipher(keyTobyt)
	if err != nil {
		fmt.Println(err)
	}

	// düz metin üzerinde şifreleme işlemi yapılıyor
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		fmt.Println(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// base64'e çeviriyoruz
	return base64.URLEncoding.EncodeToString(ciphertext)
}

// base64'e çevrilen dizeyi AES kullanarak çözüyoruz
func decrypt(this js.Value, inputs []js.Value) interface{} {

	keyTobyt := []byte(inputs[0].String())

	// base64'e çevrilen dizeyi alıyoruz
	ciphertext, _ := base64.URLEncoding.DecodeString(inputs[1].String())

	block, err := aes.NewCipher(keyTobyt)
	if err != nil {
		fmt.Println(err)
	}

	//eğer metin şifrelenmiş ise
	if len(ciphertext) < aes.BlockSize {
		return "out from chat : " + inputs[1].String()
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// dize şifresi çözülüyor
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("stranger : %s", ciphertext)
}
