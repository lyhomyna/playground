const usernameErr = document.querySelector(".username-err");
const passwordErr = document.querySelector(".password-err");
const submitBtn = document.querySelector("#submitBtn")

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
	    submitBtn.disable = true;
	    break;
	// statusNotFound - user with username not found
	case 404:
	    usernameErr.style.visibility = "collapse";
	    break;
    }
})
