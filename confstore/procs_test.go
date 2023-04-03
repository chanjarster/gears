package confstore

import (
	"crypto/rsa"
	rsautil "github.com/chanjarster/gears/util/rsa"
	"github.com/stretchr/testify/assert"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func Test_stringProcessor_Validate(t *testing.T) {
	type fields struct {
		minLength int
		maxLength int
	}
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantOk  bool
		wantErr string
	}{
		{
			fields:  fields{0, 1},
			args:    args{""},
			wantOk:  true,
			wantErr: "",
		},
		{
			fields:  fields{0, 1},
			args:    args{"a"},
			wantOk:  true,
			wantErr: "",
		},
		{
			fields:  fields{0, 1},
			args:    args{"aa"},
			wantOk:  false,
			wantErr: "len(value) not >= 0 or <= 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := StringProcessor(
				tt.fields.minLength,
				tt.fields.maxLength,
			)
			gotOk, gotErr := s.Validate(tt.args.value)
			if gotOk != tt.wantOk {
				t.Errorf("Validate() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
			if gotErr != tt.wantErr {
				t.Errorf("Validate() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func Test_stringProcessor_Convert(t *testing.T) {
	type fields struct {
		minLength int
		maxLength int
	}
	type args struct {
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		{
			fields: fields{0, 1},
			args:   args{"aa"},
			want:   "aa",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := StringProcessor(
				tt.fields.minLength,
				tt.fields.maxLength,
			)
			if got := s.Convert(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_boolProcessor_Validate(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantOk  bool
		wantErr string
	}{
		{
			args:    args{"true"},
			wantOk:  true,
			wantErr: "",
		},
		{
			args:    args{"false"},
			wantOk:  true,
			wantErr: "",
		},
		{
			args:    args{"1"},
			wantOk:  false,
			wantErr: "not bool value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &boolProcessor{}
			gotOk, gotErr := b.Validate(tt.args.value)
			if gotOk != tt.wantOk {
				t.Errorf("Validate() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
			if gotErr != tt.wantErr {
				t.Errorf("Validate() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func Test_boolProcessor_Convert(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			args: args{"true"},
			want: true,
		},
		{
			args: args{"false"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &boolProcessor{}
			if got := b.Convert(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_intProcessor_Validate(t *testing.T) {
	type fields struct {
		min        int
		minInclude bool
		max        int
		maxInclude bool
	}
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantOk  bool
		wantErr string
	}{
		{
			fields:  fields{0, true, 0, true},
			args:    args{"0"},
			wantOk:  true,
			wantErr: "",
		},
		{
			fields:  fields{0, true, 0, true},
			args:    args{"-1"},
			wantOk:  false,
			wantErr: "not >= 0",
		},
		{
			fields:  fields{0, true, 0, true},
			args:    args{"1"},
			wantOk:  false,
			wantErr: "not <= 0",
		},
		{
			fields:  fields{0, false, 0, true},
			args:    args{"0"},
			wantOk:  false,
			wantErr: "not > 0",
		},
		{
			fields:  fields{0, true, 0, false},
			args:    args{"0"},
			wantOk:  false,
			wantErr: "not < 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := IntProcessor(
				tt.fields.min,
				tt.fields.minInclude,
				tt.fields.max,
				tt.fields.maxInclude,
			)
			gotOk, gotErr := i.Validate(tt.args.value)
			if gotOk != tt.wantOk {
				t.Errorf("Validate() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
			if gotErr != tt.wantErr {
				t.Errorf("Validate() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func Test_intProcessor_Convert(t *testing.T) {
	type fields struct {
		min        int
		minInclude bool
		max        int
		maxInclude bool
	}
	type args struct {
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		{
			fields: fields{0, true, 0, true},
			args:   args{"0"},
			want:   0,
		},
		{
			fields: fields{0, true, 0, true},
			args:   args{"-1"},
			want:   -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := IntProcessor(
				tt.fields.min,
				tt.fields.minInclude,
				tt.fields.max,
				tt.fields.maxInclude,
			)
			if got := i.Convert(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_durationProcessor_Validate(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantOk  bool
		wantErr string
	}{
		{
			args:    args{"1"},
			wantOk:  false,
			wantErr: "not duration value",
		},
		{
			args:    args{"1s"},
			wantOk:  true,
			wantErr: "",
		},
		{
			args:    args{"-1s"},
			wantOk:  true,
			wantErr: "",
		},
		{
			args:    args{"1m"},
			wantOk:  true,
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &durationProcessor{}
			gotOk, gotErr := d.Validate(tt.args.value)
			if gotOk != tt.wantOk {
				t.Errorf("Validate() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
			if gotErr != tt.wantErr {
				t.Errorf("Validate() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func Test_durationProcessor_Convert(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			args: args{"1s"},
			want: time.Second,
		},
		{
			args: args{"-1s"},
			want: -time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &durationProcessor{}
			if got := d.Convert(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_urlProcessor_Validate(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantOk  bool
		wantErr string
	}{
		{
			args:    args{"abc"},
			wantOk:  false,
			wantErr: "not url value",
		},
		{
			args:    args{"http://abc"},
			wantOk:  true,
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &urlProcessor{}
			gotOk, gotErr := u.Validate(tt.args.value)
			if gotOk != tt.wantOk {
				t.Errorf("Validate() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
			if gotErr != tt.wantErr {
				t.Errorf("Validate() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func Test_urlProcessor_Convert(t *testing.T) {

	var null *url.URL

	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			args: args{"http://abc"},
			want: mustParseRequestURI("http://abc"),
		},
		{
			args: args{"http://abc.com"},
			want: mustParseRequestURI("http://abc.com"),
		},
		{
			want: null,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &urlProcessor{}
			if got := u.Convert(tt.args.value).(*url.URL); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mustParseRequestURI(rawurl string) *url.URL {
	u, _ := url.ParseRequestURI(rawurl)
	return u
}

func Test_urlStringProcessor_Validate(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantOk  bool
		wantErr string
	}{
		{
			args:    args{"abc"},
			wantOk:  false,
			wantErr: "not url value",
		},
		{
			args:    args{"http://abc"},
			wantOk:  true,
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &urlStringProcessor{}
			gotOk, gotErr := u.Validate(tt.args.value)
			if gotOk != tt.wantOk {
				t.Errorf("Validate() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
			if gotErr != tt.wantErr {
				t.Errorf("Validate() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func Test_urlStringProcessor_Convert(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			args: args{"http://abc"},
			want: "http://abc",
		},
		{
			args: args{"http://abc.com"},
			want: "http://abc.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &urlStringProcessor{}
			if got := u.Convert(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cssColorHex_Validate(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantOk  bool
		wantErr string
	}{
		{args: args{"#000000"}, wantOk: true, wantErr: ""},
		{args: args{"#ffffff"}, wantOk: true, wantErr: ""},
		{args: args{"#FFFFFF"}, wantOk: true, wantErr: ""},
		{args: args{"#00000x"}, wantOk: false, wantErr: "not css color hex value"},
		{args: args{"#000"}, wantOk: false, wantErr: "not css color hex value"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cssColorHex{}
			gotOk, gotErr := c.Validate(tt.args.value)
			if gotOk != tt.wantOk {
				t.Errorf("Validate() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
			if gotErr != tt.wantErr {
				t.Errorf("Validate() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func Test_cssColorHex_Convert(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{args: args{"#000000"}, want: "#000000"},
		{args: args{"#FFFFFF"}, want: "#FFFFFF"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cssColorHex{}
			if got := c.Convert(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_rsaPrivateKeyPkcs1_Validate(t *testing.T) {
	type args struct {
		value string
	}
	pk := "MIIEowIBAAKCAQEAnMaoRsWGQdzGwcgmCd+zyJFZg1kYhzdcvbnhX133HNyLDuA09toiZ3we/buAwyX71jTRH2jmkObrbFWaBzVBxgLTwgtMdQT8k4yODVkER1e4OALnKjPLUjs3tYa/51eIxVkQ1c0zwYZ3fHOuUu05pMaJJzXs8kWEWSvGiccFXswulIomZUeFmJhfcA8ld4Zmf19W7TBsGQRpAzksxfkxglZ1Q1YG9p9G5/GhNTnqzz+teosjW8UWFBsJBctXfyiwvFphpo0wV8vvCrJSaRmmOwJuJzvrEK0ECfYhTJGfUfeM/DKq/v3bt9f5Jn8zcYiv691Wp8Axh+1vOKRDQ4cmewIDAQABAoIBAGC0JU5qPrNXpH4ZoGUMyM/Z4FYl2fxmCN26z6rMwkXaQChv4hu2V+xvYopuYzF8t4fc0nXGdhpBZkxPzQ/PKQAj9EzIPWQdLFMzKthO5VXAZRCfscmVY0AY6Pce5DamBuZ8VucaiWdBM2jfzlK2o4zhLe6y07JpcQaz+uY3jjd3US2FrA7ZgyGHJ+cYIO4Ca0P2NyMz7/Vqp7oXMXZU5aFj6RNZcl3WhIHeWmLn0w2mL72K08TCno6lA3b5fA/L5WiOCtLJg3B5kxm21L+fJ9SZZbZg23mKevDPBFO28Nch86AnnII3IS+SS40BhLruNEvLSofgmRm90OP1ZL7gLSECgYEA1uywAoQ7DinotFix5NXFKpb2d950VdYrGoeJGJZSaiRVLrggZ1ynFv/CYNWJN8wWqMSX2161b+URr6IxnPBVh064s+jprx5+7ul0FvH9GNiJNOtWt/PskSrsHxfcNpZ9BNKAy6s1Xq5qkvtz9FdW48iyGkl3BkstY2FrUrXVmNECgYEAur0EWbDyt8QjVsaidmornZTLt9FBmTtqrFKr25ZcCLMzr/Z9GkOrAHdWI5im06rKwOObCGkh/jaQAEGU29DHxJosSwcJ/5MA7AWLUuZYpZo+iuvLp4fJnyxPp6s0ZuBa+J5ZJzYFjWCLH8P0Bq5rIER7SR1wIL3pqsyztSHFnYsCgYAM9QH7dNtWlEI6NSqsYBLd6XH8cgXTcvDwTYF/YPig+4XDJkiV0qkkItEmI1l9aqlnDUsWZ5RKpaW2T/HrwzV0zdMmzYDhFNaUMrDT/UzD5bIH5bQ5rNPzQIIxsa+N/u6UjthC7xRtm9hC+jPYZDeRzfSmzw0E7R1UY3gj6WByMQKBgDwXEOxPvXsW+Yw19ReHyKs3s2peQu6tKQF0xOqvcsQ17j8EYXNvLPrEbAqeV6irP/1cAIpvUkn8gtAiSmWFxJLmTbSC+luskVJh4lh12TbI8LFBeVJapq2605MPy5gDQCLaybONdbHtzDcXllIEFGLYxpdbirZuhl+46WczC+VLAoGBAJLBDFmIQm8LkCFupQdbptRX0nSEN3vicnodllqGT0zzWSAAMhVhuPMaXW8fIgMBcX66t47CKiqVAOZWvxUN0h6lZ/4G1GGXyzu4RTrmswABqe6kw/2KGNs7Jbs2F5SjmyKwmbJOg/IOt1qokeyibC2rNU6qLssYLBSSzpiv+2t7"
	tests := []struct {
		name    string
		args    args
		wantOk  bool
		wantErr string
	}{
		{
			args:    args{"abc"},
			wantOk:  false,
			wantErr: "invalid RSA PRIVATE KEY",
		},
		{
			args:    args{pk},
			wantOk:  true,
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &rsaPrivateKeyPkcs1{}
			gotOk, gotErr := r.Validate(tt.args.value)
			if gotOk != tt.wantOk {
				t.Errorf("Validate() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
			if gotErr != tt.wantErr {
				t.Errorf("Validate() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func Test_rsaPrivateKeyPkcs1_Convert(t *testing.T) {

	var null *rsa.PrivateKey

	type args struct {
		value string
	}
	pk := "MIIEowIBAAKCAQEAnMaoRsWGQdzGwcgmCd+zyJFZg1kYhzdcvbnhX133HNyLDuA09toiZ3we/buAwyX71jTRH2jmkObrbFWaBzVBxgLTwgtMdQT8k4yODVkER1e4OALnKjPLUjs3tYa/51eIxVkQ1c0zwYZ3fHOuUu05pMaJJzXs8kWEWSvGiccFXswulIomZUeFmJhfcA8ld4Zmf19W7TBsGQRpAzksxfkxglZ1Q1YG9p9G5/GhNTnqzz+teosjW8UWFBsJBctXfyiwvFphpo0wV8vvCrJSaRmmOwJuJzvrEK0ECfYhTJGfUfeM/DKq/v3bt9f5Jn8zcYiv691Wp8Axh+1vOKRDQ4cmewIDAQABAoIBAGC0JU5qPrNXpH4ZoGUMyM/Z4FYl2fxmCN26z6rMwkXaQChv4hu2V+xvYopuYzF8t4fc0nXGdhpBZkxPzQ/PKQAj9EzIPWQdLFMzKthO5VXAZRCfscmVY0AY6Pce5DamBuZ8VucaiWdBM2jfzlK2o4zhLe6y07JpcQaz+uY3jjd3US2FrA7ZgyGHJ+cYIO4Ca0P2NyMz7/Vqp7oXMXZU5aFj6RNZcl3WhIHeWmLn0w2mL72K08TCno6lA3b5fA/L5WiOCtLJg3B5kxm21L+fJ9SZZbZg23mKevDPBFO28Nch86AnnII3IS+SS40BhLruNEvLSofgmRm90OP1ZL7gLSECgYEA1uywAoQ7DinotFix5NXFKpb2d950VdYrGoeJGJZSaiRVLrggZ1ynFv/CYNWJN8wWqMSX2161b+URr6IxnPBVh064s+jprx5+7ul0FvH9GNiJNOtWt/PskSrsHxfcNpZ9BNKAy6s1Xq5qkvtz9FdW48iyGkl3BkstY2FrUrXVmNECgYEAur0EWbDyt8QjVsaidmornZTLt9FBmTtqrFKr25ZcCLMzr/Z9GkOrAHdWI5im06rKwOObCGkh/jaQAEGU29DHxJosSwcJ/5MA7AWLUuZYpZo+iuvLp4fJnyxPp6s0ZuBa+J5ZJzYFjWCLH8P0Bq5rIER7SR1wIL3pqsyztSHFnYsCgYAM9QH7dNtWlEI6NSqsYBLd6XH8cgXTcvDwTYF/YPig+4XDJkiV0qkkItEmI1l9aqlnDUsWZ5RKpaW2T/HrwzV0zdMmzYDhFNaUMrDT/UzD5bIH5bQ5rNPzQIIxsa+N/u6UjthC7xRtm9hC+jPYZDeRzfSmzw0E7R1UY3gj6WByMQKBgDwXEOxPvXsW+Yw19ReHyKs3s2peQu6tKQF0xOqvcsQ17j8EYXNvLPrEbAqeV6irP/1cAIpvUkn8gtAiSmWFxJLmTbSC+luskVJh4lh12TbI8LFBeVJapq2605MPy5gDQCLaybONdbHtzDcXllIEFGLYxpdbirZuhl+46WczC+VLAoGBAJLBDFmIQm8LkCFupQdbptRX0nSEN3vicnodllqGT0zzWSAAMhVhuPMaXW8fIgMBcX66t47CKiqVAOZWvxUN0h6lZ/4G1GGXyzu4RTrmswABqe6kw/2KGNs7Jbs2F5SjmyKwmbJOg/IOt1qokeyibC2rNU6qLssYLBSSzpiv+2t7"
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			args: args{pk},
			want: rsautil.MustReadPrivateKey(pk),
		},
		{
			want: null,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &rsaPrivateKeyPkcs1{}
			if got := r.Convert(tt.args.value).(*rsa.PrivateKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_rsaPublicKey_Validate(t *testing.T) {
	type args struct {
		value string
	}
	pk := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0MZHrOBQztAOvNIEQfjd4A11BV4qu7IOQZlVt8mH2HJDHnpIrweLNscrfycJhTXryw6reTwnUVmMpXfzhRjABSpjtS5RP5hxFysaWiS1lJCG6mjo2EXeDjig4Jr/ydhlR9GwFBZFBT8SqkJO6k95S5KCmorcRk5x/6eekyVaH1ZDoAiMS94H04veeSr+pAOeEzaQajkaxWLRwsHMiMcpCWO2wpSmuedZdtlYRRHQ76SGfRCjYjWT/M+1fnRUHOAC6k+xFMW8WRwAOg8emWauPSGQ5K4MHiXikkuNJ2RrT8zIqJMErFGgQ9/0OTY4rS24fichNQynS7QfktX88vidvwIDAQAB"
	tests := []struct {
		name    string
		args    args
		wantOk  bool
		wantErr string
	}{
		{
			args:    args{"abc"},
			wantOk:  false,
			wantErr: "invalid PUBLIC KEY",
		},
		{
			args:    args{pk},
			wantOk:  true,
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &rsaPublicKey{}
			gotOk, gotErr := r.Validate(tt.args.value)
			if gotOk != tt.wantOk {
				t.Errorf("Validate() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
			if gotErr != tt.wantErr {
				t.Errorf("Validate() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func Test_rsaPublicKey_Convert(t *testing.T) {
	var null *rsa.PublicKey

	type args struct {
		value string
	}
	pk := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0MZHrOBQztAOvNIEQfjd4A11BV4qu7IOQZlVt8mH2HJDHnpIrweLNscrfycJhTXryw6reTwnUVmMpXfzhRjABSpjtS5RP5hxFysaWiS1lJCG6mjo2EXeDjig4Jr/ydhlR9GwFBZFBT8SqkJO6k95S5KCmorcRk5x/6eekyVaH1ZDoAiMS94H04veeSr+pAOeEzaQajkaxWLRwsHMiMcpCWO2wpSmuedZdtlYRRHQ76SGfRCjYjWT/M+1fnRUHOAC6k+xFMW8WRwAOg8emWauPSGQ5K4MHiXikkuNJ2RrT8zIqJMErFGgQ9/0OTY4rS24fichNQynS7QfktX88vidvwIDAQAB"
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			args: args{pk},
			want: rsautil.MustReadPublicKey(pk),
		},
		{
			want: null,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &rsaPublicKey{}
			if got := r.Convert(tt.args.value).(*rsa.PublicKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_enumProcessor_Validate(t *testing.T) {
	type fields struct {
		enums []string
	}
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantOk  bool
		wantErr []string
	}{
		{
			fields: fields{[]string{"foo", "bar"}},
			args:   args{""},
			wantOk: false,
			wantErr: []string{
				"only allowed value: [foo,bar]",
				"only allowed value: [bar,foo]",
			},
		},
		{
			fields: fields{[]string{"foo", "bar"}},
			args:   args{"loo"},
			wantOk: false,
			wantErr: []string{
				"only allowed value: [foo,bar]",
				"only allowed value: [bar,foo]",
			},
		},
		{
			fields:  fields{[]string{"foo", "bar"}},
			args:    args{"foo"},
			wantOk:  true,
			wantErr: []string{""},
		},
		{
			fields:  fields{[]string{"foo", "bar"}},
			args:    args{"bar"},
			wantOk:  true,
			wantErr: []string{""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := EnumProcessor(tt.fields.enums)
			gotOk, gotErr := s.Validate(tt.args.value)
			if gotOk != tt.wantOk {
				t.Errorf("Validate() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
			assert.Contains(t, tt.wantErr, gotErr, "Validate() got unexpected error")
		})
	}
}
