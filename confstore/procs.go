package configstore

import (
	"fmt"
	rsautil "github.com/chanjarster/gears/util/rsa"
	"math"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

var (
	String             ValueProcessor = StringProcessor(0, math.MaxInt32)
	Bool               ValueProcessor = &boolProcessor{}
	Int                ValueProcessor = IntProcessor(math.MinInt32, true, math.MaxInt32, true)
	IntGtZero          ValueProcessor = IntProcessor(0, false, math.MaxInt32, true)
	IntGtEqZero        ValueProcessor = IntProcessor(0, true, math.MaxInt32, true)
	Duration           ValueProcessor = &durationProcessor{}
	Url                ValueProcessor = &urlProcessor{}
	UrlString          ValueProcessor = &urlStringProcessor{}
	CssColorHex        ValueProcessor = &cssColorHex{}
	RsaPrivateKeyPkcs1 ValueProcessor = &rsaPrivateKeyPkcs1{}
	RsaPublicKey       ValueProcessor = &rsaPublicKey{}
)

// TODO write examples

type stringProcessor struct {
	minLength, maxLength int
}

func StringProcessor(minLength, maxLength int) *stringProcessor {
	return &stringProcessor{minLength: minLength, maxLength: maxLength}
}

func (s *stringProcessor) Validate(value string) (ok bool, err string) {
	ok =  len(value) >= s.minLength && len(value) <= s.maxLength
	if !ok {
		err = fmt.Sprintf("len(value) not >= %d or <= %d", s.minLength, s.maxLength)
	}
	return ok, err
}

func (s *stringProcessor) Convert(value string) interface{} {
	return value
}

type boolProcessor struct{}

func (b *boolProcessor) Validate(value string) (ok bool, err string) {
	if value != "true" && value != "false" {
		return false, "not bool value"
	}
	return true, ""
}

func (b *boolProcessor) Convert(value string) interface{} {
	v, _ := strconv.ParseBool(value)
	return v
}

func IntProcessor(min int, minInclude bool, max int, maxInclude bool) *intProcessor {
	return &intProcessor{min: min, minInclude: minInclude, max: max, maxInclude: maxInclude}
}

type intProcessor struct {
	min        int
	minInclude bool
	max        int
	maxInclude bool
}

func (i *intProcessor) Validate(value string) (ok bool, err string) {
	v, e := strconv.Atoi(value)
	if e != nil {
		return false, "not int value"
	}
	if i.minInclude {
		if v < i.min {
			return false, "not >= " + strconv.Itoa(i.min)
		}
	} else {
		if v <= i.min {
			return false, "not > " + strconv.Itoa(i.min)
		}
	}
	if i.maxInclude {
		if v > i.max {
			return false, "not <= " + strconv.Itoa(i.max)
		}
	} else {
		if v >= i.max {
			return false, "not < " + strconv.Itoa(i.max)
		}
	}
	return true, ""
}

func (i *intProcessor) Convert(value string) interface{} {
	v, _ := strconv.Atoi(value)
	return v
}

type durationProcessor struct{}

func (d *durationProcessor) Validate(value string) (ok bool, err string) {
	_, e := time.ParseDuration(value)
	if e != nil {
		return false, "not duration value"
	}
	return true, ""
}

func (d *durationProcessor) Convert(value string) interface{} {
	v, _ := time.ParseDuration(value)
	return v
}

type urlProcessor struct{}

func (u *urlProcessor) Validate(value string) (ok bool, err string) {
	if value == "" {
		return true, ""
	}
	_, e := url.ParseRequestURI(value)
	if e != nil {
		return false, "not url value"
	}
	return true, ""
}

func (u *urlProcessor) Convert(value string) interface{} {
	if value == "" {
		return nil
	}
	v, _ := url.ParseRequestURI(value)
	return v
}

type urlStringProcessor struct{}

func (u *urlStringProcessor) Validate(value string) (ok bool, err string) {
	if value == "" {
		return true, ""
	}
	_, e := url.ParseRequestURI(value)
	if e != nil {
		return false, "not url value"
	}
	return true, ""
}

func (u *urlStringProcessor) Convert(value string) interface{} {
	return value
}

type cssColorHex struct{}

func (c *cssColorHex) Validate(value string) (ok bool, err string) {
	if value == "" {
		return true, ""
	}

	match, error := regexp.MatchString("#[a-fA-F0-9]{6}", value)
	if error != nil {
		return false, error.Error()
	}
	if !match {
		return false, "not css color hex value"
	}
	return true, ""
}

func (c *cssColorHex) Convert(value string) interface{} {
	return value
}

type rsaPrivateKeyPkcs1 struct{}

func (r *rsaPrivateKeyPkcs1) Validate(value string) (ok bool, err string) {
	if value == "" {
		return true, ""
	}
	_, error := rsautil.ReadPrivateKey(value)
	if error != nil {
		return false, error.Error()
	}
	return true, ""
}

func (r *rsaPrivateKeyPkcs1) Convert(value string) interface{} {
	if value == "" {
		return nil
	}
	v, _ := rsautil.ReadPrivateKey(value)
	return v
}

type rsaPublicKey struct{}

func (r *rsaPublicKey) Validate(value string) (ok bool, err string) {
	if value == "" {
		return true, ""
	}
	_, error := rsautil.ReadPublicKey(value)
	if error != nil {
		return false, error.Error()
	}
	return true, ""
}

func (r *rsaPublicKey) Convert(value string) interface{} {
	if value == "" {
		return nil
	}
	v, _ := rsautil.ReadPublicKey(value)
	return v
}
