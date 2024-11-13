const usernameErr = document.querySelector(".username-err");
const passwordErr = document.querySelector(".password-err");
const submitBtn = document.querySelector("#submit-btn")

const username = document.querySelector(".username");
username.addEventListener("input", async (_) => {
    const opt = {
	method: "post",
	body: JSON.stringify(username.value),
    };
    const response = await fetch("/users", opt);
    switch (response.status) {
	// statusOK - user with username found
	case 200:
	    usernameErr.style.visibility = "visible";
	    break;
	// statusNotFound - user with username not found
	case 404:
	    usernameErr.style.visibility = "collapse";
	    break;
    }
    enableButton()
})

const password = document.querySelector(".password")
const passwordRepeat = document.querySelector(".password-repeat")

password.addEventListener("change", (_) => {
    if (passwordRepeat.value.trim() === "" && password.value.trim() === "") {
	passwordErr.style.visibility = "collapse";
    } else if (passwordRepeat.value.trim()) {
	passwordErr.style.visibility = password.value.trim() === passwordRepeat.value.trim() 
	    ? "collapse"
	    : "visible"
    }
    enableButton()
})
passwordRepeat.addEventListener("change", (_) => {
    if (passwordRepeat.value.trim() === "" && password.value.trim() === "") {
	passwordErr.style.visibility = "collapse";
    } else if (password.value.trim()) {
	passwordErr.style.visibility = (password.value.trim() === passwordRepeat.value.trim()) 
	    ? "collapse" 
	    : "visible";
    }
    enableButton()
})

const firstname = document.querySelector("#firstname")
const lastname = document.querySelector("#lastname")

firstname.addEventListener("change", (_) => enableButton())
lastname.addEventListener("change", (_) => enableButton())

function enableButton() {
    console.log("enable button")
    if (usernameErr.style.visibility !== "visible" 
	&& passwordErr.style.visibility !== "visible"
	&& username.value.trim() !== ""
	&& password.value.trim() !== ""
	&& passwordRepeat.value.trim() !== ""
	&& firstname.value.trim() !== ""
	&& lastname.value.trim() !== "") {
	console.log("nu vse blyat', vse ok")
        submitBtn.disabled = false;
    } else {
        submitBtn.disabled = true;
    }
}

