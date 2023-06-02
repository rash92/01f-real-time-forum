function showCommentEdit(elem) {
    console.log("Hello got here")
    var x = document.getElementById("commentEditor")
    var y = document.getElementById("editbutton")
    var z = document.getElementById("finalEditButton")

    if (elem.value == "showComment") {
        x.style.display = "block";
        y.style.display = "none";
        z.style.display = "block";
    }
}