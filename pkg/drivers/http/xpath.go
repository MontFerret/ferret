package http

import (
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/PuerkitoBio/goquery"
)

type XPathNavigator struct {
	selection *goquery.Selection
}

func NewXPathNavigator(selection *goquery.Selection) *XPathNavigator {

}

func (xn *XPathNavigator) Query(expression string) (drivers.HTMLElement, error) {}

func (xn *XPathNavigator) QueryAll(expression string) (drivers.HTMLElement, error) {}
