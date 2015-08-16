
package command

import (
	"os"
	"errors"
	"strconv"
	"strings"
	"io"
	"io/ioutil"
	"net/url"
	"net/http"
)

const URL_META = "http://www.youtube.com/get_video_info?&video_id="
var FORMATS []string = []string{"3gp", "mp4", "flv", "webm", "avi"}

type Video struct {
	Id, Title, Author, Keywords, Thumbnail_url string
	Avg_rating float32
	View_count,	Length_seconds int
	Formats []Format
}

type Format struct {
	Itag int
	Video_type, Quality, Url string
}

func Get(video_id string) (Video, error) {
	query_string, err := fetchMeta(video_id)
	if err != nil {
		return Video{}, err
	}

	meta, err := parseMeta(video_id, query_string)

	if err != nil {
		return Video{}, err
	}

	return meta, nil
}

func (video *Video) Download(index int, filename string) error {
	out, err := os.Create(filename)
	defer out.Close()

	if err != nil {
		return errors.New("Unable to write to file " + filename)
	}

	resp, err := http.Get(video.Formats[index].Url)
	defer resp.Body.Close()

	if err != nil {
		return errors.New("Unable to download video content from Yotutube")
	}

	io.Copy(out, resp.Body)

	return nil
}

func (video *Video) GetExtension(index int) string {
	for i := 0; i<len(FORMATS); i++ {
		if strings.Contains(video.Formats[index].Video_type, FORMATS[i]) {
			return FORMATS[i]
		}
	}

	return "avi"
}

func fetchMeta(video_id string) (string, error) {
	resp, err := http.Get(URL_META + video_id)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	query_string, _ := ioutil.ReadAll(resp.Body)

	return string(query_string), nil
}

func parseMeta(video_id, query_string string) (Video, error) {
	u, _ := url.Parse("?" + query_string)

	query := u.Query()

	if query.Get("errorcode") != "" || query.Get("status") == "fail" {
		return Video{}, errors.New(query.Get("reason"))
	}

	video := Video{
		Id: video_id,
		Title: query.Get("title"),
		Author: query.Get("author"),
		Keywords: query.Get("keywords"),
		Thumbnail_url: query.Get("thumbnail_url"),
	}
	
	v, _ := strconv.Atoi(query.Get("view_count"))
	video.View_count = v

	r, _ := strconv.ParseFloat(query.Get("avg_rating"), 32)
	video.Avg_rating = float32(r)

	l, _ := strconv.Atoi(query.Get("length_seconds"))
	video.Length_seconds = l

	format_params := strings.Split(query.Get("url_encoded_fmt_stream_map"), ",")

	for _, f := range format_params {
		furl, _ := url.Parse("?" + f)
		fquery := furl.Query()

		itag, _ := strconv.Atoi(fquery.Get("itag"))

		video.Formats = append(
			video.Formats,
			Format{
				Itag: itag,
				Video_type: fquery.Get("type"),
				Quality: fquery.Get("quality"),
				Url: fquery.Get("url") + "&signature=" + fquery.Get("sig"),
			})
	}

	return video, nil
}