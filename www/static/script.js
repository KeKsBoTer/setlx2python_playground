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
    console.log("transpiling...")
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
            codeEditorPython.setValue(res.code)
            errorBox.style.display="none"
        })
        .catch(e => {
            errorBox.innerText = e
            errorBox.style.display="initial"
        })
}