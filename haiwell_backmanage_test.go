package tests

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// func TestHmiProjectDownload(t *testing.T) {
// 	const (
// 		HAIWELL_PROJECT_FILEPATH = "./assets/3.40.0.14.hwdev"
// 	)
// 	password := os.Getenv("HAIWELL_PROJECT_ENCRYPT_PASSWORD")
//
// 	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?_key=%s", HAIWELL_PROJECT_FILEPATH, password))
// 	if err != nil {
// 		t.Fatalf("Failed to open the database: %v", err)
// 	}
// 	defer db.Close()
//
// 	rows, err := db.Query("SELECT * FROM haiwell")
// 	if err != nil {
// 		t.Fatalf("Failed to query the comand: %v", err)
// 	}
//
// 	defer rows.Close()
// }

func TestHmiProjectUpload(t *testing.T) {

	const (
		BACKMANAGE_PROTOCOL     = "http"
		BACKMANAGE_IP           = "192.168.22.23"
		BACKMANAGE_UPDATE_PORT  = 81
		BACKMANAGE_UONLINE_PORT = 82

		BACKMANAGE_GET_HMI_INFO                             = "/update/getHmiInfo"
		BACKMANAGE_GET_HMI_INFO_UPLAOD_PASSWORD_STATE_FIELD = "loadPwdState"
		BACKMANAGE_GET_HMI_INFO_UPLOAD_PROJECT_PREMIT_FIELD = "uploadPrjPermit"

		BACKMANAGE_PROJECT_FILE_COUNT                    = "/uonline/getFileCount/"
		BACKMANAGE_PROJECT_FILE_COUNT_FILE_CUTSIZE_FILED = "cutsize"
		BACKMANAGE_PROJECT_FILE_COUNT_FILE_TYPE_FIELD    = "fileType"
		BACKMANAGE_PROJECT_FILE_COUNT_FILE_COUNT_FIELD   = "fileCount"
		BACKMANAGE_PROJECT_FILE_COUNT_FILE_MD5_FIELD     = "md5"
		BACKMANAGE_PROJECT_FILE_COUNT_FILE_TYPE_VALUE    = "project"
		BACKMANAGE_PROJECT_FILE_COUNT_FILE_CUTSIZE_VALUE = 512 * 1024

		BACKMANAGE_UPLOAD_PROJECT_FILE                = "/uonline/uploadProject/"
		BACKMANAGE_UPLOAD_PROJECT_FILE_STORAGE_DIR    = "./assets/"
		BACKMANAGE_UPLOAD_PROJECT_FILE_CUTSIZE_FILED  = "cutSize"
		BACKMANAGE_UPLOAD_PROJECT_FILE_CUTSIZE_VALUE  = 512 * 1024
		BACKMANAGE_UPLOAD_PROJECT_FILE_INDEX_FILED    = "index"
		BACKMANAGE_UPLOAD_PROJECT_FILE_PASSWORD_FILED = "password"
		BACKMANAGE_UPLOAD_PROJECT_FILE_PASSWORD_VALUE = "b51b74011735cf017faeda4520932bab"
	)

	requestUrl := fmt.Sprintf("%s://%s:%d%s", BACKMANAGE_PROTOCOL, BACKMANAGE_IP, BACKMANAGE_UPDATE_PORT, BACKMANAGE_GET_HMI_INFO)
	m, err := httpPostForm(requestUrl, nil)
	if err != nil {
		t.Fatalf("Failed to get the HMI info: %v", err)
	}

	if m[BACKMANAGE_GET_HMI_INFO_UPLAOD_PASSWORD_STATE_FIELD].(float64) != 1 ||
		m[BACKMANAGE_GET_HMI_INFO_UPLOAD_PROJECT_PREMIT_FIELD].(float64) != 1 {
		t.Fatalf("The HMI upload project permission is not allowed")
	}

	requestUrl = fmt.Sprintf("%s://%s:%d%s", BACKMANAGE_PROTOCOL, BACKMANAGE_IP, BACKMANAGE_UONLINE_PORT, BACKMANAGE_PROJECT_FILE_COUNT)

	values := url.Values{
		BACKMANAGE_PROJECT_FILE_COUNT_FILE_CUTSIZE_FILED: {fmt.Sprintf("%d", BACKMANAGE_PROJECT_FILE_COUNT_FILE_CUTSIZE_VALUE)},
		BACKMANAGE_PROJECT_FILE_COUNT_FILE_TYPE_FIELD:    {BACKMANAGE_PROJECT_FILE_COUNT_FILE_TYPE_VALUE},
	}

	m, err = httpPostForm(requestUrl, values)
	if err != nil {
		t.Fatalf("Failed to get the project file count: %v", err)
	}

	requestUrl = fmt.Sprintf("%s://%s:%d%s", BACKMANAGE_PROTOCOL, BACKMANAGE_IP, BACKMANAGE_UONLINE_PORT, BACKMANAGE_UPLOAD_PROJECT_FILE)

	values.Set(BACKMANAGE_UPLOAD_PROJECT_FILE_CUTSIZE_FILED, fmt.Sprintf("%d", BACKMANAGE_UPLOAD_PROJECT_FILE_CUTSIZE_VALUE))
	// passwordMd5 := md5.New()
	// passwordMd5.Write([]byte(BACKMANAGE_UPLOAD_PROJECT_FILE_PASSWORD_VALUE))
	// values.Set(BACKMANAGE_UPLOAD_PROJECT_FILE_PASSWORD_FILED, string(passwordMd5.Sum(nil)))
	values.Set(BACKMANAGE_UPLOAD_PROJECT_FILE_PASSWORD_FILED, BACKMANAGE_UPLOAD_PROJECT_FILE_PASSWORD_VALUE)

	f, err := os.Create(fmt.Sprintf("%s%s.hwdev", BACKMANAGE_UPLOAD_PROJECT_FILE_STORAGE_DIR, m[BACKMANAGE_PROJECT_FILE_COUNT_FILE_MD5_FIELD].(string)))
	if err != nil {
		t.Fatalf("Failed to create the project file: %v", err)
	}

	for i := range int(m[BACKMANAGE_PROJECT_FILE_COUNT_FILE_COUNT_FIELD].(float64)) {
		values.Set(BACKMANAGE_UPLOAD_PROJECT_FILE_INDEX_FILED, fmt.Sprintf("%d", i))

		resp, err := http.PostForm(requestUrl, values)
		if err != nil {
			t.Fatalf("Failed to upload the project file: %v", err)
		}

		buff, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read the response body: %v", err)
		}

		_, err = f.Write(buff)
		if err != nil {
			t.Fatalf("Failed to write the project file: %v", err)
		}

		resp.Body.Close()
	}
}

func httpPostForm(url string, values url.Values) (map[string]any, error) {
	resp, err := http.PostForm(url, values)
	if err != nil {
		return nil, err
	}

	buff, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	m := make(map[string]any)

	if err := json.Unmarshal(buff, &m); err != nil {
		return nil, err
	}

	return m, nil
}

func calculateSign(vs url.Values) string {

	keys := make([]string, 0, len(vs))

	for key := range vs {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	sb := strings.Builder{}

	for _, key := range keys {
		sb.WriteString(vs.Get(key))
	}

	h := md5.New()
	h.Write([]byte(sb.String()))
	return string(h.Sum(nil))

}
