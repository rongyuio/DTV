package bilibili

import "testing"

func TestLoginUID(t *testing.T) {
	cases := []struct {
		name   string
		cookie string
		want   string
	}{
		{"未登录", "", "0"},
		{"只有游客 cookie", "buvid3=abc; b_nut=1", "0"},
		{"普通 uid", "SESSDATA=x; DedeUserID=12345678; bili_jct=y", "12345678"},
		// 2^53 之后 float64 不再精确，这条钉住「不转数字」的行为。
		{"超出 float64 精确范围", "DedeUserID=9007199254740993", "9007199254740993"},
		{"19 位 uid", "DedeUserID=1234567890123456789; SESSDATA=x", "1234567890123456789"},
		{"前后空白", " DedeUserID = 42 ", "0"},
		{"值带空白", "DedeUserID= 42 ; bili_jct=x", "42"},
		// DedeUserID__ckMd5 是另一个字段，不能被误当成 uid。
		{"md5 字段不匹配", "DedeUserID__ckMd5=deadbeef", "0"},
		{"非数字值", `DedeUserID=12abc`, "0"},
		{"空值", "DedeUserID=", "0"},
		{"注入尝试", `DedeUserID=1,"x":2`, "0"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client := NewClient()
			client.SetCookie(c.cookie)
			if got := client.loginUID(); got != c.want {
				t.Fatalf("loginUID() = %q, want %q", got, c.want)
			}
		})
	}
}

func TestAllDigits(t *testing.T) {
	yes := []string{"0", "42", "1234567890123456789"}
	no := []string{"", " 1", "1 ", "-1", "1.0", "1e5", "١٢٣"}

	for _, s := range yes {
		if !allDigits(s) {
			t.Errorf("allDigits(%q) = false, want true", s)
		}
	}
	for _, s := range no {
		if allDigits(s) {
			t.Errorf("allDigits(%q) = true, want false", s)
		}
	}
}
