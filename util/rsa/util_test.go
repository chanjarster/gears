package rsautil

import (
	"testing"
)

func TestReadPrivateKey(t *testing.T) {

	type args struct {
		pemStr string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args:    args{"MIIEowIBAAKCAQEAnMaoRsWGQdzGwcgmCd+zyJFZg1kYhzdcvbnhX133HNyLDuA09toiZ3we/buAwyX71jTRH2jmkObrbFWaBzVBxgLTwgtMdQT8k4yODVkER1e4OALnKjPLUjs3tYa/51eIxVkQ1c0zwYZ3fHOuUu05pMaJJzXs8kWEWSvGiccFXswulIomZUeFmJhfcA8ld4Zmf19W7TBsGQRpAzksxfkxglZ1Q1YG9p9G5/GhNTnqzz+teosjW8UWFBsJBctXfyiwvFphpo0wV8vvCrJSaRmmOwJuJzvrEK0ECfYhTJGfUfeM/DKq/v3bt9f5Jn8zcYiv691Wp8Axh+1vOKRDQ4cmewIDAQABAoIBAGC0JU5qPrNXpH4ZoGUMyM/Z4FYl2fxmCN26z6rMwkXaQChv4hu2V+xvYopuYzF8t4fc0nXGdhpBZkxPzQ/PKQAj9EzIPWQdLFMzKthO5VXAZRCfscmVY0AY6Pce5DamBuZ8VucaiWdBM2jfzlK2o4zhLe6y07JpcQaz+uY3jjd3US2FrA7ZgyGHJ+cYIO4Ca0P2NyMz7/Vqp7oXMXZU5aFj6RNZcl3WhIHeWmLn0w2mL72K08TCno6lA3b5fA/L5WiOCtLJg3B5kxm21L+fJ9SZZbZg23mKevDPBFO28Nch86AnnII3IS+SS40BhLruNEvLSofgmRm90OP1ZL7gLSECgYEA1uywAoQ7DinotFix5NXFKpb2d950VdYrGoeJGJZSaiRVLrggZ1ynFv/CYNWJN8wWqMSX2161b+URr6IxnPBVh064s+jprx5+7ul0FvH9GNiJNOtWt/PskSrsHxfcNpZ9BNKAy6s1Xq5qkvtz9FdW48iyGkl3BkstY2FrUrXVmNECgYEAur0EWbDyt8QjVsaidmornZTLt9FBmTtqrFKr25ZcCLMzr/Z9GkOrAHdWI5im06rKwOObCGkh/jaQAEGU29DHxJosSwcJ/5MA7AWLUuZYpZo+iuvLp4fJnyxPp6s0ZuBa+J5ZJzYFjWCLH8P0Bq5rIER7SR1wIL3pqsyztSHFnYsCgYAM9QH7dNtWlEI6NSqsYBLd6XH8cgXTcvDwTYF/YPig+4XDJkiV0qkkItEmI1l9aqlnDUsWZ5RKpaW2T/HrwzV0zdMmzYDhFNaUMrDT/UzD5bIH5bQ5rNPzQIIxsa+N/u6UjthC7xRtm9hC+jPYZDeRzfSmzw0E7R1UY3gj6WByMQKBgDwXEOxPvXsW+Yw19ReHyKs3s2peQu6tKQF0xOqvcsQ17j8EYXNvLPrEbAqeV6irP/1cAIpvUkn8gtAiSmWFxJLmTbSC+luskVJh4lh12TbI8LFBeVJapq2605MPy5gDQCLaybONdbHtzDcXllIEFGLYxpdbirZuhl+46WczC+VLAoGBAJLBDFmIQm8LkCFupQdbptRX0nSEN3vicnodllqGT0zzWSAAMhVhuPMaXW8fIgMBcX66t47CKiqVAOZWvxUN0h6lZ/4G1GGXyzu4RTrmswABqe6kw/2KGNs7Jbs2F5SjmyKwmbJOg/IOt1qokeyibC2rNU6qLssYLBSSzpiv+2t7"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ReadPrivateKey(tt.args.pemStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildRsaPrivateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}

}

func TestReadPublicKey(t *testing.T) {
	type args struct {
		pemStr string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args:    args{"MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0MZHrOBQztAOvNIEQfjd4A11BV4qu7IOQZlVt8mH2HJDHnpIrweLNscrfycJhTXryw6reTwnUVmMpXfzhRjABSpjtS5RP5hxFysaWiS1lJCG6mjo2EXeDjig4Jr/ydhlR9GwFBZFBT8SqkJO6k95S5KCmorcRk5x/6eekyVaH1ZDoAiMS94H04veeSr+pAOeEzaQajkaxWLRwsHMiMcpCWO2wpSmuedZdtlYRRHQ76SGfRCjYjWT/M+1fnRUHOAC6k+xFMW8WRwAOg8emWauPSGQ5K4MHiXikkuNJ2RrT8zIqJMErFGgQ9/0OTY4rS24fichNQynS7QfktX88vidvwIDAQAB"},
			wantErr: false,
		},
		{
			args:    args{"MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnMaoRsWGQdzGwcgmCd+zyJFZg1kYhzdcvbnhX133HNyLDuA09toiZ3we/buAwyX71jTRH2jmkObrbFWaBzVBxgLTwgtMdQT8k4yODVkER1e4OALnKjPLUjs3tYa/51eIxVkQ1c0zwYZ3fHOuUu05pMaJJzXs8kWEWSvGiccFXswulIomZUeFmJhfcA8ld4Zmf19W7TBsGQRpAzksxfkxglZ1Q1YG9p9G5/GhNTnqzz+teosjW8UWFBsJBctXfyiwvFphpo0wV8vvCrJSaRmmOwJuJzvrEK0ECfYhTJGfUfeM/DKq/v3bt9f5Jn8zcYiv691Wp8Axh+1vOKRDQ4cmewIDAQAB"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ReadPublicKey(tt.args.pemStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildRsaPublicKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}

}