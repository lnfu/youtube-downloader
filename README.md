一個輸入表單可以輸入網址以及選擇要

1. only 影片
2. only 音樂
3. 影片＋音樂

按下按鈕或是 ENTER 後送出

（如果網址無效會有紅字）

會發出
POST /api/download

```
POST
Request Body:
{
    "url": "bHLpxwhELEs",
    "type": "a" # type = a | v | av
}
```

---

```
bestVideoAndAudioCombinedFormat := "bv+ba/b"
// bestVideoAndAudioFormat := "bv,ba"
bestVideoFormat:= "bv"
bestAudioFormat:= "ba"

output, err := exec.Command("yt-dlp", "-f", format, "https://www.youtube.com/watch?v=" + video_id).Output()
if err != nil {
    // log.Fatal(err)
} else {
    fmt.Printf("%s\n", output)
}

```

```
sudo -u postgres psql
```
