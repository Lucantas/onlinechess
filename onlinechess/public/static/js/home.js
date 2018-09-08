window.onload = function(){
    var signin = document.getElementById("signin");
    var signup = document.getElementById("signup");

    signin.onclick = function(){
        document.querySelectorAll("[data-transition=move-left]").forEach(function(item){
           //while (parseInt(item.style.left.replace("px","")) != -900){
                console.log("left", item.style.left)
                item.style.left = `${parseInt(item.style.left.replace("px","")) - 10}px`     
           //}
        })
    }
    signup.onclick = function(){
        document.querySelectorAll("[data-transition=move-left]").forEach(function(item){
          
        })
    }
    
}

