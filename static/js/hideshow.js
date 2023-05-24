function hideShow(elem) {
    var x = document.getElementById("user-posts");
    var a = document.getElementById("posts")
    var y = document.getElementById("user-comments");
    var b = document.getElementById("comments")
    var z = document.getElementById("liked-user-posts");
    var c = document.getElementById("likes")
    var d = document.getElementById("disliked-user-posts");
    var e = document.getElementById("dislikes")
    if (elem.value == "posts") {
        x.style.display = "block";
        y.style.display = "none";
        z.style.display = "none";
        d.style.display = "none";
        a.style.textDecoration = "underline";
        b.style.textDecoration = "none";
        c.style.textDecoration = "none";
        e.style.textDecoration = "none";
    } else if (elem.value == "comments") {
        x.style.display = "none";
        y.style.display = "block";
        z.style.display = "none";
        d.style.display = "none";
        a.style.textDecoration = "none";
        b.style.textDecoration = "underline";
        c.style.textDecoration = "none";
        e.style.textDecoration = "none";
    } else if (elem.value == "liked-posts") {
        x.style.display = "none";
        y.style.display = "none";
        z.style.display = "block";
        d.style.display = "none";
        a.style.textDecoration = "none"
        b.style.textDecoration = "none"
        c.style.textDecoration = "underline"
        e.style.textDecoration = "none";
    } else if (elem.value == "disliked-posts") {
        x.style.display = "none";
        y.style.display = "none";
        z.style.display = "none";
        d.style.display = "block";
        a.style.textDecoration = "none"
        b.style.textDecoration = "none"
        c.style.textDecoration = "none"
        e.style.textDecoration = "underline";
    }
}