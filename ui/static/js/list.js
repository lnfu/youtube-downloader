function saveMediaListToSessionStorage(list) {
  sessionStorage.setItem('userMediaList', JSON.stringify(list));
}

function getMediaListFromSessionStorage() {
  const storedMediaList = sessionStorage.getItem('userMediaList');
  return storedMediaList ? JSON.parse(storedMediaList) : [];
}


document.getElementById('form').addEventListener('submit', function(event) {
    event.preventDefault();
    const url = this.elements['url'].value;
    const type = this.elements['type'].value;
    
    const v = getYouTubeVideoId(url)
    if (v) {
        if (type === "video") {
            list.push({
                "v": v,
                "type": "v"
            })      
        }

        if (type === "audio") {
            list.push({
                "v": v,
                "type": "a"
            })      
        }
    }
});



function getOriginInfo(item) {
    fetch('/api/get_origin_info?'+"v="+item.v)
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }
            return response.json();
        })
        .then(jsonData => {

            item.title = jsonData.title
            item.duration = secondsToHMS(jsonData.duration)
            
        })
        .catch(error => {
            console.error('Fetch error:', error);
        });
}

function getMediaInfo(item) {
    fetch('/api/get_media_info?'+"v="+item.v+'&t='+item.type)
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }
            return response.json();
        })
        .then(jsonData => {

            item.filename = jsonData.filename
            item.key = jsonData.key
            item.status = jsonData.status
        })
        .catch(error => {
            console.error('Fetch error:', error);
        });
}



function renderItem(item, index) {
    const table = document.getElementById("list")
        var row = table.rows[index];
        if (!row) {
            row = table.insertRow(index);
            for (var colIndex = 0; colIndex < 6; colIndex++) {
                var cell = row.cells[colIndex];
                if (!cell) {
                    cell = row.insertCell(colIndex);
                }

                if (colIndex !== 1) {
                    cell.style.textAlign = 'center';
                }
            }

            const downloadButton = document.createElement('button');
            downloadButton.textContent = '下載';
            downloadButton.addEventListener('click', function () {
                downloadFile(item.key, item.filename);
            });

            row.cells[5].appendChild(downloadButton);
            row.cells[5].style.textAlign = 'center';
        }

    row.cells[0].textContent = index;
    row.cells[1].textContent = item.title;

    switch (item.type) {
        case "v":
            row.cells[2].textContent = "影片";
            break;
        case "a":
            row.cells[2].textContent = "音樂";
            break;
    
        default:
            break;
    }


    row.cells[3].textContent = item.duration;

    switch (item.status) {
        case "running":
            row.cells[4].textContent = "處理中";
            break;
        case "done":
            row.cells[4].textContent = "完成";
            break;
        case "failure":
            row.cells[4].textContent = "失敗";
            break;
    
        default:
            break;
    }


}

const list = getMediaListFromSessionStorage(); 

function loop() {

    var index = 0
    list.forEach(item => {

        index++

        if (!item.title) {
            getOriginInfo(item)
        }

        if (!item.status || item.status === "running") {
            getMediaInfo(item)
        }    
        
        renderItem(item, index);

    });

    setTimeout(loop, 1000);

}

loop();



function getYouTubeVideoId(url) {
    const regex = /^https:\/\/www\.youtube\.com\/watch\?v=([a-zA-Z0-9_-]+)$/;
    const match = regex.exec(url);
    return match ? match[1] : null;
}



function secondsToHMS(seconds) {
    var hours = Math.floor(seconds / 3600);
    var minutes = Math.floor((seconds % 3600) / 60);
    var remainingSeconds = seconds % 60;

    var result = pad(hours) + ':' + pad(minutes) + ':' + pad(remainingSeconds);
    return result;
}

function pad(number) {
    return (number < 10) ? '0' + number : number;
}


function downloadFile(key, filename) {
    fetch('/api/get_file', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            key: key,
            filename: filename,
        }),
    })
    .then(response => {
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        return response.blob();
    })
    .then(blob => {
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = filename; 
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
    })
    .catch(error => {
        console.error('Fetch error:', error);
    });
}
