function clear_mutations() {
    var last;
    const mutations = get_mutations();
    while (last = mutations.lastChild) {
        mutations.removeChild(last);
    }
}

function clear_gif() {
    var last;
    const mutations = get_gif();
    while (last = mutations.lastChild) {
        mutations.removeChild(last);
    }
}

function get_mutations() {
    return document.getElementById("mutations");
}

function get_gif() {
    return document.getElementById("gif");
}

function image_clicked() {
    const xhr = new XMLHttpRequest();
    const image = this.image;
    xhr.open('GET', '/choose_image?id=' + image.id);
    xhr.onload = function() {
        if (xhr.status === 200) {
            get_images.bind({
                id: image.id
            })();
        } else {
            alert('Request failed.  Returned status of ' + xhr.status);
        }
    }
    xhr.send();
}

function history_clicked() {
    const xhr = new XMLHttpRequest();
    const image = this.image;
    xhr.open('GET', '/get_image?id=' + image.id);
    xhr.onload = function() {
        if (xhr.status === 200) {
            const json = JSON.parse(xhr.responseText);
            draw_images(json.images);
            draw_gif(json.gif);
        } else {
            alert('Request failed.  Returned status of ' + xhr.status);
        }
    }
    xhr.send();
}

function create_image(image) {
    const div = document.createElement("div");
    div.setAttribute("class", "container");
    const img_node = document.createElement("img");
    img_node.setAttribute("class", "mutant clickable");
    img_node.setAttribute("src", "data:image/png;base64," + image.bytes);
    img_node.onclick = image_clicked.bind({
        image: image
    });
    div.appendChild(img_node);
    const a = document.createElement("span");
    a.setAttribute("class", "clickable top-left");
    a.innerText = "ID: " + image.id;
    a.onclick = history_clicked.bind({
        image: image
    });
    div.appendChild(a);
    return div;
}

function draw_images(images) {
    clear_mutations();
    const mutations = get_mutations();
    for (const img_i in images) {
        mutations.appendChild(create_image(images[img_i]));
    }
}

function draw_gif(gif) {
    clear_gif();
    const el = get_gif();
    const im = document.createElement("img");
    im.setAttribute("src", "data:iamge/gif;base64," + gif);
    im.setAttribute("class", "mutant");
    el.appendChild(im);
}

function get_images() {
    const xhr = new XMLHttpRequest();
    if (this.id != undefined) {
        xhr.open('GET', '/mutate_image?id=' + this.id);
    } else {
        xhr.open('GET', '/get_images');
    }
    xhr.onload = function() {
        if (xhr.status === 200) {
            draw_images(JSON.parse(xhr.responseText).images);
        } else {
            alert('Request failed.  Returned status of ' + xhr.status);
        }
    };
    xhr.send();
}