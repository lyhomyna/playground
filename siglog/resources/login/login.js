const username = document.querySelector("#username")
const password = document.querySelector("#password")
const submitBtn = document.querySelector("#submit-btn")

submitBtn.addEventListener("click", async (event) => {
    event.preventDefault()

    const opt = {
	method: "post",
	body: JSON.stringify({
	    username: username.value.trim(),
	    password: password.value.trim(),
	}),
    }

    const response = await fetch("/login", opt)
    if (response.status == 200) {
	location.reload()
    } else if (response.status == 401) {
	console.log("Unauthorized.")
    }
})
username.addEventListener("change", () => enableButton())
password.addEventListener("change", () => enableButton())

function enableButton() {
    console.log("enable button")
    if (username.value.trim() !== "" && password.value.trim() !== "") {
	console.log("nu vse blyat', vse ok")
        submitBtn.disabled = false;
    } else {
        submitBtn.disabled = true;
    }
}
