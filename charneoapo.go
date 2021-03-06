package charneoapo

import (
	"io"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	gq "github.com/PuerkitoBio/goquery"
	"github.com/ikeikeikeike/go-bracmeister"
	"github.com/ikeikeikeike/gopkg/convert"
	behavior "github.com/ikeikeikeike/gopkg/net/http"
	"github.com/ikeikeikeike/gopkg/str"
)

const EndPoint = "http://neoapo.com"

func tee(r io.Reader, debug bool) io.Reader {
	if !debug {
		return r
	}
	return io.TeeReader(r, os.Stdout)
}

type Neoapo struct {
	*behavior.UserBehavior
	doc *gq.Document

	Unit  string
	Debug bool
}

func NewNeoapo() *Neoapo {
	return &Neoapo{
		UserBehavior: behavior.NewUserBehavior(),
		Unit:         "cm",
		Debug:        false,
	}
}

func (w *Neoapo) Doc(url string) (*gq.Document, error) {
	resp, err := w.Behave(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return gq.NewDocumentFromResponse(resp)
}

func (w *Neoapo) Do(path, page string) error {
	doc, err := w.Doc(EndPoint + "/" + path + "/" + url.QueryEscape(page))

	if err != nil {
		return err
	}

	w.doc = doc
	return nil
}

func (w *Neoapo) Name() (r string) {
	w.doc.Find(`dl dt:contains(名前)`).Each(func(i int, s *gq.Selection) {
		r = str.Clean(s.Next().Text())
	})
	return
}

func (w *Neoapo) Kana() (r string) {
	w.doc.Find(`dl dt:contains(名前)`).Each(func(i int, s *gq.Selection) {
		r, _ = s.Next().Attr("title")
		r = str.Clean(r)
	})
	return
}

func (w *Neoapo) Product() (r string) {
	var ok bool
	w.doc.Find(`.profile_related a`).Each(func(i int, s *gq.Selection) {
		if r != "" {
			return
		}
		r, ok = s.ChildrenFiltered("img").Attr("title")
		if !ok {
			return
		}
		r = str.Clean(r)
	})
	return
}

func (w *Neoapo) Birthday() (r time.Time) {
	w.doc.Find(`dl dt:contains(誕生日)`).Each(func(i int, s *gq.Selection) {
		r, _ = time.Parse("2006年1月2日", str.MustClean(s.Next().Text()))
	})
	return
}

func (w *Neoapo) Blood() (r string) {
	w.doc.Find(`dl dt:contains(血液型)`).Each(func(i int, s *gq.Selection) {
		r = str.Clean(strings.Replace(s.Next().Text(), "型", "", -1))
	})
	return
}

var reNum = regexp.MustCompile(`(\d+)`)

func (w *Neoapo) Height() (r int) {
	w.doc.Find(`dl dt:contains(身長)`).Each(func(i int, s *gq.Selection) {
		text := s.Next().Text()
		if strings.Contains(text, w.Unit) {
			r, _ = convert.StrTo(str.Clean(reNum.FindString(text))).Int()
		}
	})
	return
}

func (w *Neoapo) Weight() (r int) {
	var err error
	w.doc.Find(`dl dt:contains(体重)`).Each(func(i int, s *gq.Selection) {
		text := s.Next().Text()
		if strings.Contains(text, "kg") {
			text = str.Clean(strings.Replace(text, "kg", "", -1))
			r, err = convert.StrTo(text).Int()
			if err != nil {
				f, _ := convert.StrTo(text).Float32()
				r = int(f)
			}
		}
	})
	return
}

func (w *Neoapo) BWH() (r string) {
	w.doc.Find(`dl dt:contains(ｽﾘｰｻｲｽﾞ)`).Each(func(i int, s *gq.Selection) {
		text := s.Next().Text()
		if strings.Contains(text, "B") {
			r = str.Clean(text)
		}
	})
	return
}

func (w *Neoapo) Bust() (r int) {
	bhw := strings.Split(w.BWH(), "/")
	if len(bhw) > 0 {
		r, _ = convert.StrTo(str.Clean(reNum.FindString(bhw[0]))).Int()
	}
	return
}

func (w *Neoapo) Waist() (r int) {
	bhw := strings.Split(w.BWH(), "/")
	if len(bhw) > 1 {
		r, _ = convert.StrTo(str.Clean(reNum.FindString(bhw[1]))).Int()
	}
	return
}

func (w *Neoapo) Hip() (r int) {
	bhw := strings.Split(w.BWH(), "/")
	if len(bhw) > 2 {
		r, _ = convert.StrTo(str.Clean(reNum.FindString(bhw[2]))).Int()
	}
	return
}

func (w *Neoapo) Bracup() (r string) {
	var re = regexp.MustCompile(`\(\w\)`)

	bhw := strings.Split(w.BWH(), "/")
	if len(bhw) > 0 {
		r = re.FindString(bhw[0])
		r = strings.Replace(strings.Replace(r, "(", "", -1), ")", "", -1)
	}

	if r == "" {
		h, b, w := w.Height(), w.Bust(), w.Waist()
		if h > 10 && b > 10 && w > 10 {
			r = bracmeister.Calc(h, b, w, true).Cup
		}
	}

	r = strings.ToUpper(str.Clean(r))
	return
}

func (w *Neoapo) Comment() (r string) {
	w.doc.Find(`dl dt:contains(コメント)`).Each(func(i int, s *gq.Selection) {
		r = str.Clean(s.Next().Next().Text())
	})
	return
}

func (w *Neoapo) AnimeName() (r string) {
	w.doc.Find(`dl dt:contains(作品名)`).Each(func(i int, s *gq.Selection) {
		r = str.Clean(s.Next().Text())
	})
	return
}

func (w *Neoapo) AnimeAlias() (r string) {
	w.doc.Find(`dl dt:contains(通名\/略称)`).Each(func(i int, s *gq.Selection) {
		r = str.Clean(s.Next().Text())
	})
	return
}

func (w *Neoapo) AnimeAuthor() (r string) {
	w.doc.Find(`dl dt:contains(原作者)`).Each(func(i int, s *gq.Selection) {
		r = str.Clean(s.Next().Text())
	})
	return
}

func (w *Neoapo) AnimeWorks() (r string) {
	w.doc.Find(`dl dt:contains(制作会社)`).Each(func(i int, s *gq.Selection) {
		r = str.Clean(s.Next().Text())
	})
	return
}

func (w *Neoapo) AnimeRelease() (r time.Time) {
	w.doc.Find(`dl dt:contains(制作年)`).Each(func(i int, s *gq.Selection) {
		r, _ = time.Parse("2006", str.MustClean(s.Next().Text()))
	})
	return
}

func (w *Neoapo) AnimeUrl() (r string) {
	w.doc.Find(`dl dt:contains(公式サイト)`).Each(func(i int, s *gq.Selection) {
		r = str.Clean(s.Next().Text())
	})
	return
}
