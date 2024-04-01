package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	var islem int
	fmt.Println("********************************")
	fmt.Println("*         Banka Sistemi        *")
	fmt.Println("********************************")
	fmt.Println("    1. Banka Hesabı Oluştur     ")
	fmt.Println("    2. Banka Hesabına Giriş Yap ")
	fmt.Println("********************************")
	fmt.Println("Yapmak istediğiniz işlemi seçin: ")
	fmt.Scan(&islem)

	if islem == 1 {
		createBankAccount()
	} else if islem == 2 {
		var id, pass, money string

		fmt.Println("********************************")
		fmt.Println("*           Giriş Yap          *")
		fmt.Println("********************************")
		fmt.Println("     Kimlik numaranızı girin    ")
		fmt.Scan(&id)
		fmt.Println("         Şifrenizi girin        ")
		fmt.Scan(&pass)
		if login("users.txt", id, pass) {
			fmt.Println("Giriş başarılı! Hoş geldiniz.\n\n")
			fmt.Println("********************************")
			fmt.Println("*           İŞLEMLER           *")
			fmt.Println("*         1. Para Çek          *")
			fmt.Println("*         2. Para Yatır        *")
			fmt.Println("********************************")
			fmt.Println("Yapmak istediğiniz işlemi seçin: ")
			fmt.Scan(&islem)

			if islem == 1 {
				fmt.Println("Çekmek istediğiniz tutarı girin:")
				fmt.Scanln(&money)
				withdrawMoney(id, money)
			} else if islem == 2 {
				fmt.Println("Yatırmak istediğiniz tutarı girin:")
				fmt.Scanln(&money)
				depositMoney(id, money)
			} else {
				fmt.Print("Hatalı İşlem Yaptınız!")
			}
		} else {
			fmt.Println("Giriş bilgileri geçersiz. Lütfen tekrar deneyin.")
		}

	} else {
		fmt.Print("Hatalı İşlem Yaptınız!")
	}
}
func createBankAccount() {
	var id string
	var pass string
	fmt.Println("********************************")
	fmt.Println("*         Hesap Oluştur        *")
	fmt.Println("********************************")
	fmt.Println("     Kimlik numaranızı girin    ")
	fmt.Scan(&id)
	fmt.Println("         Şifrenizi girin        ")
	fmt.Scan(&pass)
	writeToFile("users.txt", id+":"+pass+"\n")
}

func login(filename, id, pass string) bool {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Dosya okuma hatası:", err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) == 2 && parts[0] == id && parts[1] == pass {
			return true
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Dosya okuma hatası:", err)
		return false
	}

	return false
}

func writeToFile(filename, data string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

func depositMoney(id string, money string) {
	writeToFile("money.txt", id+":"+money+"\n")
}

func withdrawMoney(id string, amount string) {
	balance, err := getBalance(id)
	if err != nil {
		fmt.Println("Hesap bilgisi alınamadı:", err)
		return
	}

	withdrawAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		fmt.Println("Geçersiz miktar:", err)
		return
	}

	if balance < withdrawAmount {
		fmt.Println("Yetersiz bakiye!")
		return
	}

	// Yeni bakiyeyi hesapla
	newBalance := balance - withdrawAmount

	// Yeni bakiyeyi dosyaya yaz
	err = updateBalance(id, newBalance)
	if err != nil {
		fmt.Println("Bakiye güncellenemedi:", err)
		return
	}

	fmt.Println("Para çekme işlemi başarılı! Yeni bakiye:", newBalance)
}

func getBalance(id string) (float64, error) {
	file, err := os.Open("money.txt")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) == 2 && parts[0] == id {
			return strconv.ParseFloat(parts[1], 64)
		}
	}

	return 0, fmt.Errorf("hesap bulunamadı")
}

func updateBalance(id string, newBalance float64) error {
	tempFile, err := os.CreateTemp("", "temp_balances.txt")
	if err != nil {
		return err
	}
	defer tempFile.Close()

	file, err := os.Open("money.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) == 2 && parts[0] == id {
			line = id + ":" + strconv.FormatFloat(newBalance, 'f', -1, 64)
		}
		tempFile.WriteString(line + "\n")
	}

	err = os.Rename(tempFile.Name(), "money.txt")

	if err != nil {
		return err
	}

	return nil
}
