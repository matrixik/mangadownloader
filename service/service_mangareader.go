package service

import (
	"code.google.com/p/go-html-transform/css/selector"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
)

const (
	serviceMangaReadeChapterDigitCount = 4
)

var (
	mangareader = &MangaReaderService{
		Hosts: []string{
			"www.mangareader.net",
			"mangareader.net",
		},
	}

	serviceMangaReaderHtmlSelectorIdentifyManga, _   = selector.Selector("#chapterlist")
	serviceMangaReaderHtmlSelectorIdentifyChapter, _ = selector.Selector("#pageMenu")
	serviceMangaReaderHtmlSelectorMangaName, _       = selector.Selector("h2.aname")
	serviceMangaReaderHtmlSelectorMangaChapters, _   = selector.Selector("#chapterlist a")
	serviceMangaReaderHtmlSelectorChapterName, _     = selector.Selector("#mangainfo h1")
	serviceMangaReaderHtmlSelectorChapterPages, _    = selector.Selector("#pageMenu option")
	serviceMangaReaderHtmlSelectorPageImage, _       = selector.Selector("#img")

	serviceMangaReaderRegexpChapterName, _ = regexp.Compile("^.* (\\d*)$")

	serviceMangaReaderFormatChapter = "%0" + strconv.Itoa(serviceMangaReadeChapterDigitCount) + "d"
)

func init() {
	mangareader.UrlBase = new(url.URL)
	mangareader.UrlBase.Scheme = "http"
	mangareader.UrlBase.Host = mangareader.Hosts[0]

	RegisterService("mangareader", mangareader)
}

type MangaReaderService Service

func init() {
	RegisterService("mangareader", mangareader)
}

func (service *MangaReaderService) Supports(u *url.URL) bool {
	return stringSliceContains(mangareader.Hosts, u.Host)
}

func (service *MangaReaderService) Identify(u *url.URL) (interface{}, error) {
	if !service.Supports(u) {
		return nil, errors.New("Not supported")
	}

	rootNode, err := HttpGetHtml(u, service.HttpRetry)
	if err != nil {
		return nil, err
	}

	identifyMangaNodes := serviceMangaReaderHtmlSelectorIdentifyManga.Find(rootNode)
	if len(identifyMangaNodes) == 1 {
		manga := &Manga{
			Url:     u,
			Service: service,
		}
		return manga, nil
	}

	identifyChapterNodes := serviceMangaReaderHtmlSelectorIdentifyChapter.Find(rootNode)
	if len(identifyChapterNodes) == 1 {
		chapter := &Chapter{
			Url:     u,
			Service: service,
		}
		return chapter, nil
	}

	return nil, errors.New("Unknown url")
}

func (service *MangaReaderService) MangaName(manga *Manga) (string, error) {
	rootNode, err := HttpGetHtml(manga.Url, service.HttpRetry)
	if err != nil {
		return "", err
	}

	nameNodes := serviceMangaReaderHtmlSelectorMangaName.Find(rootNode)
	if len(nameNodes) != 1 {
		return "", errors.New("Name node not found")
	}
	nameNode := nameNodes[0]
	if nameNode.FirstChild == nil {
		return "", errors.New("Name text node not found")
	}
	nameTextNode := nameNode.FirstChild
	name := nameTextNode.Data

	return name, nil
}

func (service *MangaReaderService) MangaChapters(manga *Manga) ([]*Chapter, error) {
	rootNode, err := HttpGetHtml(manga.Url, service.HttpRetry)
	if err != nil {
		return nil, err
	}

	linkNodes := serviceMangaReaderHtmlSelectorMangaChapters.Find(rootNode)

	chapters := make([]*Chapter, 0, len(linkNodes))
	for _, linkNode := range linkNodes {
		chapterUrl := urlCopy(mangareader.UrlBase)
		chapterUrl.Path = htmlGetNodeAttribute(linkNode, "href")
		chapter := &Chapter{
			Url:     chapterUrl,
			Service: service,
		}
		chapters = append(chapters, chapter)
	}

	return chapters, nil
}

func (service *MangaReaderService) ChapterName(chapter *Chapter) (string, error) {
	rootNode, err := HttpGetHtml(chapter.Url, service.HttpRetry)
	if err != nil {
		return "", err
	}

	nameNodes := serviceMangaReaderHtmlSelectorChapterName.Find(rootNode)
	if len(nameNodes) != 1 {
		return "", errors.New("Name node not found")
	}
	nameNode := nameNodes[0]
	if nameNode.FirstChild == nil {
		return "", errors.New("Name text node not found")
	}
	nameTextNode := nameNode.FirstChild
	name := nameTextNode.Data
	matches := serviceMangaReaderRegexpChapterName.FindStringSubmatch(name)
	if matches == nil {
		return "", errors.New("Invalid name format")
	}
	name = matches[1]
	nameInt, err := strconv.Atoi(name)
	if err != nil {
		return "", err
	}
	name = fmt.Sprintf(serviceMangaReaderFormatChapter, nameInt)

	return name, nil
}

func (service *MangaReaderService) ChapterPages(chapter *Chapter) ([]*Page, error) {
	rootNode, err := HttpGetHtml(chapter.Url, service.HttpRetry)
	if err != nil {
		return nil, err
	}

	optionNodes := serviceMangaReaderHtmlSelectorChapterPages.Find(rootNode)

	pages := make([]*Page, 0, len(optionNodes))
	for _, optionNode := range optionNodes {
		pageUrl := urlCopy(mangareader.UrlBase)
		pageUrl.Path = htmlGetNodeAttribute(optionNode, "value")
		page := &Page{
			Url:     pageUrl,
			Service: service,
		}
		pages = append(pages, page)
	}

	return pages, nil
}

func (service *MangaReaderService) PageImageUrl(page *Page) (*url.URL, error) {
	rootNode, err := HttpGetHtml(page.Url, service.HttpRetry)
	if err != nil {
		return nil, err
	}

	imgNodes := serviceMangaReaderHtmlSelectorPageImage.Find(rootNode)
	if len(imgNodes) != 1 {
		return nil, errors.New("Image node not found")
	}
	imgNode := imgNodes[0]

	imageUrl, err := url.Parse(htmlGetNodeAttribute(imgNode, "src"))
	if err != nil {
		return nil, err
	}

	return imageUrl, nil
}

func (service *MangaReaderService) String() string {
	return "MangaReaderService"
}
