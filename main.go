package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Lütfen bir link girin!")
		return
	}
	link := os.Args[1]

	dosyaAdi := link
	dosyaAdi = strings.ReplaceAll(dosyaAdi, "https://", "")
	dosyaAdi = strings.ReplaceAll(dosyaAdi, "http://", "")
	dosyaAdi = strings.ReplaceAll(dosyaAdi, "www.", "")
	dosyaAdi = strings.ReplaceAll(dosyaAdi, "/", "")

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var siteHtml string
	var ekranGoruntusu []byte
	var linkler []string

	fmt.Printf("'%s' sitesine bağlanılıyor...\n", link)

	err := chromedp.Run(ctx,
		chromedp.Navigate(link),
		chromedp.Sleep(30*time.Second),
		chromedp.OuterHTML("html", &siteHtml),
		chromedp.Evaluate(`Array.from(document.querySelectorAll('a')).map(a => a.href)`, &linkler),
		chromedp.FullScreenshot(&ekranGoruntusu, 90),
	)

	if err != nil {
		fmt.Println("Hata oluştu:", err)
		return
	}

	htmlDosyaIsmi := dosyaAdi + "_data.html"
	os.WriteFile(htmlDosyaIsmi, []byte(siteHtml), 0644)
	fmt.Println("HTML kaydedildi:", htmlDosyaIsmi)

	resimDosyaIsmi := dosyaAdi + "_screenshot.png"
	os.WriteFile(resimDosyaIsmi, ekranGoruntusu, 0644)
	fmt.Println("Resim kaydedildi:", resimDosyaIsmi)

	if len(linkler) > 0 {
		linkDosyaIsmi := dosyaAdi + "_linkler.txt"
		linkMetni := ""
		for _, l := range linkler {
			linkMetni += l + "\n"
		}
		os.WriteFile(linkDosyaIsmi, []byte(linkMetni), 0644)
		fmt.Printf("Linkler kaydedildi (%d adet): %s\n", len(linkler), linkDosyaIsmi)
	} else {
		fmt.Println("Link bulunamadığı için link dosyası oluşturulmadı.")
	}

	fmt.Println("İşlem tamam!")
}
