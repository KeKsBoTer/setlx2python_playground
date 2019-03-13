let codeEditorSetlx, codeEditorPython

window.onload = () => {
    let codeSetlx = document.getElementById("code-setlx")
    let codePython = document.getElementById("code-python")
    codeEditorSetlx = CodeMirror.fromTextArea(codeSetlx, {
        lineNumbers: true,
        autofocus: true,
        indentUnit: 4,
        mode: "setlx",
        readOnly: false
    })
    codeEditorPython = CodeMirror.fromTextArea(codePython, {
        lineNumbers: true,
        autofocus: true,
        indentUnit: 4,
        mode: "python",
        readOnly: true
    })

    let changeTimer = undefined
    codeEditorSetlx.on("change", function (e) {
        if (changeTimer !== undefined) {
            clearInterval(changeTimer)
            changeTimer = undefined
        }
        changeTimer = setTimeout(transpileCode, 1000)
    })
}


function transpileCode() {
    let errorBox = document.getElementById("transpile-error")
    let code = codeEditorSetlx.getValue()
    fetch('/transpile', {
            method: 'post',
            headers: {
                'Accept': 'application/json, text/plain, */*',
                'Content-Type': 'application/json'
            },
            body: code
        })
        .then(res => {
            if (!res.ok) {
                throw Error(res.statusText);
            }
            return res.json()
        })
        .then(res => {
            if (res.success === true) {
                codeEditorPython.setValue(res.code)
                errorBox.style.display = "none"
            } else {
                errorBox.innerText = res.error
                errorBox.style.display = ""
            }
        })
        .catch(e => {
            errorBox.innerText = e
            errorBox.style.display = ""
        })
}

let executingSetlX = false

function runSetlXCode() {
    if (executingSetlX)
        return
    let code = codeEditorSetlx.getValue()

    // check if code is empty
    if (!code)
        return
    let out = document.getElementById("setlX-output")
    executingSetlX = true
    out.innerHTML = "<span class='info'>" + "Waiting for remote server..." + "</span>"
    fetch("/run/setlx", {
        method: "POST",
        headers: {
            "Content-Type": "text/plain; charset=utf-8",
        },
        body: code
    }).then(async (response) => {
        executingSetlX = false
        if (response.ok)
            return response.json()
        let errorMsg = await response.text()
        throw errorMsg ? errorMsg : response.statusText
    }).then(async (json) => {
        out.innerHTML = ""
        let log = ""
        for (var msg of json.events) {
            log += `<span class="${msg.Kind}">${msg.Text}</span>`
        }
        out.innerHTML += log;
        out.innerHTML += `<span class="info">Program exited.</span>`
        out.scrollY = out.scrollHeight
        out.scrollTop = out.scrollHeight
    }).catch((e) => {
        executingSetlX = false
        out.innerHTML = `<span class="stderr">${e}</span>`
    })
}

let executingPython = false

function runPythonCode() {
    if (executingPython)
        return

    let code = codeEditorPython.getValue()

    // check if code is empty
    if (!code)
        return
    let out = document.getElementById("python-output")
    executingPython = true
    out.innerHTML = "<span class='info'>" + "Waiting for remote server..." + "</span>"

    fetch("/run/python", {
        method: "POST",
        headers: {
            "Content-Type": "text/plain; charset=utf-8",
        },
        body: code
    }).then(async (response) => {
        executingPython = false
        if (response.ok)
            return response.json()
        let errorMsg = await response.text()
        throw errorMsg ? errorMsg : response.statusText
    }).then(async (json) => {
        out.innerHTML = ""
        let log = ""
        for (var msg of json.events) {
            log += `<span class="${msg.Kind}">${msg.Text}</span>`
        }
        out.innerHTML += log;
        out.innerHTML += `<span class="info">Program exited.</span>`
        out.scrollY = out.scrollHeight
        out.scrollTop = out.scrollHeight
    }).catch((e) => {
        executingPython = false
        out.innerHTML = `<span class="stderr">${e}</span>`
    })
}