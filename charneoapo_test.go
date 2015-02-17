package charneoapo

import "testing"

func TestSimple(t *testing.T) {
	c := NewNeoapo()

	err := c.Do("2100101")

	if err != nil {
		t.Fatal(err)
	}

	if c.Name() != "本田未央" {
		t.Errorf("Unexpected Neoapo.Name: %s", c.Name())
	}

	if c.Kana() != "ほんだみお" {
		t.Errorf("Unexpected Neoapo.Kana: %s", c.Kana())
	}

	if c.Product() != "THE IDOLM@STER CINDERELLA GIRLS アイドルマスターシンデレラガールズ" {
		t.Errorf("Unexpected Neoapo.Product: %s", c.Product())
	}

	if c.Birthday().Unix() != 849398400 {
		t.Errorf("Unexpected Neoapo.Birthday: %s", c.Birthday())
	}

	if c.Blood() != "B" {
		t.Errorf("Unexpected Neoapo.Blood: %s", c.Blood())
	}

	if c.Height() != 161 {
		t.Errorf("Unexpected Neoapo.Height: %d", c.Height())
	}
	if c.Weight() != 46 {
		t.Errorf("Unexpected Neoapo.Weight: %d", c.Weight())
	}

	if c.BWH() != "B84(D)/W58/H87" {
		t.Errorf("Unexpected Neoapo.BWH: %s", c.BWH())
	}
	if c.Bust() != 84 {
		t.Errorf("Unexpected Neoapo.Bust: %d", c.Bust())
	}
	if c.Waist() != 58 {
		t.Errorf("Unexpected Neoapo.Waist: %d", c.Waist())
	}
	if c.Hip() != 87 {
		t.Errorf("Unexpected Neoapo.Hip: %d", c.Hip())
	}

	if c.Bracup() != "D" {
		t.Errorf("Unexpected Neoapo.Bracup: %s", c.Bracup())
	}

	if c.Comment() != "パッションを選ぶとチュートリアルで仲間になる。覚醒美希っぽい髪型とピンクのジャージが特徴的で、明るく元気一杯。口癖は「えへへっ」" {
		t.Errorf("Unexpected Neoapo.Comment: %s", c.Comment())
	}
}
