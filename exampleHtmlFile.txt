package main

func GetHtmlFile() (file string) {
	return `
<!DOCTYPE html>
<html lang="de">
    <head>
        <base href="/">
        <link rel="shortcut icon" href="favicon.ico" type="image/x-icon" />
        <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
        <meta name="viewport" content="width=900, user-scalable=yes, initial-scale=1" />
        
        <title>Raum E008 (APB/E008/U) - Andreas-Pfitzmann-Bau</title>
        <link href="style.css" rel="stylesheet" type="text/css" />
        <script src="cn.min.js" type="text/javascript" charset="utf-8"></script>
        <script>
            if (typeof cn === 'undefined') {
                document.write(unescape("%3Cscript src='cn.js' type='text/javascript' charset='utf-8'%3E%3C/script%3E"));
            }
        </script>
        <script src="kinetic-v5.1.0.min.js" type="text/javascript" charset="utf-8"></script>
        <script type="text/javascript">
            var mobile = false;
            /**
             * I18n strings in the preferred language with fallback to "de".
             */
            let i18n = {};

            var RAplan = new function () {

                var png_file_name;
                var quali_steps;
                var subpics_size;
                var canvas_width;
                var canvas_height;
                var data_canv_width;
                var data_canv_height;
                var g_tilex;
                var g_tiley;
                var g_xtrans;
                var g_ytrans;
                var scale = 1;
                var factor = 12 / 11;
                var layx = 0;
                var layy = 0;
                var quali_cur = -1;
                var subpics = [];
                var visibleX;
                var visibleY;
                var cont_div;
                var canvas;
                var cache;
                var derraum;
                var index_names = 0;
                var lasteditbarf;       // to reset background-color of previous selected div

                this.init = function () {

                    canvas_width = 200;
                    canvas_height = 200;

                    cont_div = document.getElementById("lageplan");
                    canvas = new Kinetic.Stage({container: "lageplan", width: canvas_width, height: canvas_height});

                    if (mobile) {
                        cont_div.addEventListener("touchstart", function () {
                            location.href = document.getElementsByTagName('base')[0].href + "etplan/apb/00/raum/542100.2130";
                        }, false);
                    }

                    png_file_name = "apb00";
                    quali_steps = [1,2,4,8];
                    subpics_size = 480;

                    data_canv_width = 737;
                    data_canv_height = 480;
                    layx = 0.5 * canvas_width;
                    layy = 0.5 * canvas_height;
                    scale = 1;

                    quali_cur = 0;

                    visibleX = [1, -1];
                    visibleY = [1, -1];

                    for (var i = 0; i < quali_steps.length; i++) {
                        var h = Math.ceil((data_canv_width * quali_steps[i]) / subpics_size);
                        var v = Math.ceil((data_canv_height * quali_steps[i]) / subpics_size);
                        subpics.push({h: h, v: v});
                    }

                    // Objektlayer zeichnen ----------------------------------------------------------------------------------
            derraum = new Kinetic.Layer({x:canvas_width*0.5, y:canvas_height*0.5, listening:false, name:"derraum"});
derraumData = {points: [[262.5583644691027, 216.30707938962388, 298.3036516173984, 216.30707938962388, 298.3036516173984, 206.19979129941612, 299.68415925898773, 206.19979129941612, 299.68415925898773, 157.18868853272267, 262.5583644691027, 157.18868853272267], [281.00691252495096, 186.59109721648872, -331.3742052101149, -180.88160914309879]], fill:"#AE0000"};
for(var i = 0; i < derraumData.points.length; i++) {
derraum.add(new Kinetic.Line({closed:true, points: derraumData.points[i], fill:derraumData.fill, listening: false}));
}
canvas.add(derraum);


                    // Cachelayer vorbereiten ----------------------------------------------------------------------------------
                    cache = new Kinetic.Layer({x: canvas_width * 0.5, y: canvas_height * 0.5, listening: false});
                    canvas.add(cache);

                    // verfügbare / gelieferte Breite ----------------------------------------------------------------------------------
                    var asp_w = canvas_width /737;
                    var asp_h = canvas_height /480;
                    var min = asp_w < asp_h ? asp_w : asp_h;

                    // Weiteres
                    scale = Math.pow(factor, Math.floor(Math.log(min) / Math.log(factor)));
                    scale *= 1.5;

                    var schwerpunkt = [0, 0];
                    var ages = 0;
                    var posA = new Array();
                    for (var k = 0; k < derraumData.points.length; k++) {
                        posA.push(k);
                    }
                    for (var k = 0; k < posA.length; k++) {
                        // Schwerpunkt bestimmen 
                        var mpX, mpY;
                        var n = derraumData.points[posA[k]].length * 0.5;
                        var ai, atmp = 0, xtmp = 0, ytmp = 0;
                        for (var i = n - 1, j = 0; j < n; i = j, j++) {
                            ai = derraumData.points[posA[k]][2 * i] * derraumData.points[posA[k]][2 * j + 1] - derraumData.points[posA[k]][2 * j] * derraumData.points[posA[k]][2 * i + 1];
                            atmp += ai;
                            xtmp += (derraumData.points[posA[k]][2 * j] + derraumData.points[posA[k]][2 * i]) * ai;
                            ytmp += (derraumData.points[posA[k]][2 * j + 1] + derraumData.points[posA[k]][2 * i + 1]) * ai;
                        }
                        if (atmp != 0) {
                            mpX = xtmp / (3 * atmp);
                            mpY = ytmp / (3 * atmp);
                            schwerpunkt[0] += (mpX * Math.abs(atmp));
                            schwerpunkt[1] += (mpY * Math.abs(atmp));
                            ages += Math.abs(atmp);
                        }
                        // Ende Schwerpunkt
                    }
                    if (schwerpunkt[0] != 0 && schwerpunkt[1] != 0) {
                        schwerpunkt[0] /= ages;
                        schwerpunkt[1] /= ages;
                        move(-schwerpunkt[0] * scale, -schwerpunkt[1] * scale);
                    }

                    // Zeichnen
                    updateGlobalTilePos();
                    rescale();
                    drawAllVisibleSubpics();
                    redrawAllLayers();

                    //initLSK
                    initMoreLSK();

                    localStorage.getItem('priv-stmt') === 'false' ? cn.Main.hidePrivStmt() : cn.Main.showPrivStmt();

                    // init after page content is loaded:
                    document.addEventListener("DOMContentLoaded", cn.loadI18n("accessibility", "de", i18n, loadAccessibilityIconsAndTable, []), false);
//                    window.addEventListener("load", cn.loadI18n("de", i18n, loadAccessibilityIconsAndTable, []), false);
                };

                function loadAccessibilityIconsAndTable() {
                    let torun = new Map();
                    torun.set(200, function (responseText) {
                        let resp = JSON.parse(responseText);
                        let noinfo = false;
                        let i18n_yes = "ja";
                        let i18n_partly = "teilweise";
                        let i18n_no = "nein";
                        let i18n_na = "keine/nicht ausreichend Informationen";
                        if (resp !== null && !resp.error && resp.type) {
                            if (resp.type === 12) {
                                cn.$("icon_wheelchair").style.display = "none";
                                cn.$("icon_barfwc").style.display = "none";
                                cn.$("icon_recreationroom").style.display = "none";
                                cn.$("icon_wheelchairlecturer").style.display = "none";
                                cn.$("icon_wheelchairtable").style.display = "none";
                                cn.$("icon_markedsteps").style.display = "none";
                            } else if (resp.type === 14) {
                                cn.$("icon_wheelchair").style.display = "none";
                                cn.$("icon_barflift").style.display = "none";
                                cn.$("icon_recreationroom").style.display = "none";
                                cn.$("icon_wheelchairlecturer").style.display = "none";
                                cn.$("icon_wheelchairtable").style.display = "none";
                                cn.$("icon_markedsteps").style.display = "none";
                            } else if (resp.type === 26) {
                                cn.$("icon_wheelchair").style.display = "none";
                                cn.$("icon_barfwc").style.display = "none";
                                cn.$("icon_barflift").style.display = "none";
                                cn.$("icon_wheelchairlecturer").style.display = "none";
                                cn.$("icon_wheelchairtable").style.display = "none";
                                cn.$("icon_markedsteps").style.display = "none";
                            } else if (resp.type === 22 || resp.type === 23) {
                                cn.$("icon_wheelchair").style.display = "none";
                                cn.$("icon_barfwc").style.display = "none";
                                cn.$("icon_barflift").style.display = "none";
                                cn.$("icon_recreationroom").style.display = "none";
                            } else {
                                cn.$("icon_wheelchair").style.display = "none";
                                cn.$("icon_barfwc").style.display = "none";
                                cn.$("icon_barflift").style.display = "none";
                                cn.$("icon_recreationroom").style.display = "none";
                                cn.$("icon_wheelchairlecturer").style.display = "none";
                                cn.$("icon_wheelchairtable").style.display = "none";
                                cn.$("icon_markedsteps").style.display = "none";
                            }
                            if (resp.accessibility) {
                                // build table from the allinfo object showing the available information as text
                                if (resp.accessibility.hasOwnProperty("allinfo")) {
                                    let table = cn.$("barf_info");

                                    let tr_pic = document.createElement("TR");
                                    for (let key in resp.accessibility.allinfo) {
                                        if (resp.accessibility.allinfo[key] !== "") {
                                            let tr = document.createElement("TR");
                                            let td_key = document.createElement("TD");
                                            td_key.className = "blue klapp";
                                            td_key.innerHTML = i18n["accessibility." + key];
                                            tr.appendChild(td_key);
                                            let td_value = document.createElement("TD");
                                            td_value.style.width = "470px";
                                            let category = resp.accessibility.allinfo[key];
                                            if (key !== "pic") {
                                                for (let key1 in category) {
                                                    let div = document.createElement("DIV");
                                                    let valuetype = cn.getJSONValueType(category[key1]);
                                                    let value = "";
                                                    if (valuetype === "boolean") {
                                                        if (category[key1]) {
                                                            value = i18n_yes;
                                                        } else {
                                                            value = i18n_no;
                                                        }
                                                    } else {
                                                        value = category[key1];
                                                    }
                                                    div.style.height = "22px";
                                                    div.innerHTML = "<b>" + i18n["accessibility." + key1] + ":</b> " + value;
                                                    let button = document.createElement("INPUT");
                                                    button.className = "editinputbutton";
                                                    button.value = "edit";
                                                    button.title = "Editieren";
                                                    button.type = "button";
                                                    button.style.display = "none";
                                                    button.style.marginLeft = "10px";
                                                    button.onclick = function () {
                                                        editBarfData(div, key1, valuetype);
                                                    };
                                                    div.appendChild(button);
                                                    //TODO: activate with touch as well
                                                    div.onmouseover = function () {
                                                        button.style.display = "inline";
                                                    };
                                                    div.onmouseout = function () {
                                                        button.style.display = "none";
                                                    };
                                                    td_value.appendChild(div);
                                                }
                                                tr.appendChild(td_value);
                                            } else {
                                                tr_pic.appendChild(td_key);
                                                let img = document.createElement("IMG");
                                                img.src = category;
                                                td_value.appendChild(img);
                                                tr_pic.appendChild(td_value);
                                            }
                                            table.appendChild(tr);
                                        }
                                    }
                                    if (tr_pic.children.length === 1) {
                                        table.appendChild(tr_pic);
                                    }

                                    if (resp.accessibilityLastUpdate !== null) {
                                        let lastUpdate = new Date(resp.accessibilityLastUpdate);
                                        cn.$("barf_info_effective").innerHTML = "Stand: " + cn.getDateString(lastUpdate);
                                    }

                                    cn.$("barf_div").style.display = "block";
                                }

                                if (resp.accessibility.hasOwnProperty("barf_wc")) {
                                    let i18n_barfwc = "Barrierefreies WC";
                                    if (resp.accessibility.barf_wc === "FULLY_BARF") {
                                        cn.accessibilityIconYes("icon_barfwc", "popup_barfwc",
                                            i18n_barfwc + ": " + i18n_yes);
                                    } else if (resp.accessibility.barf_wc === "SEMI_BARF") {
                                        cn.accessibilityIconPartly("icon_barfwc", "popup_barfwc",
                                            i18n_barfwc + ": " + i18n_partly);
                                    } else if (resp.accessibility.barf_wc === "NOT_BARF") {
                                        cn.accessibilityIconNo("icon_barfwc", "popup_barfwc",
                                            i18n_barfwc + ": " + i18n_no);
                                    } else if (resp.accessibility.barf_wc === "NO_INFO") {
                                        cn.accessibilityIconNa("icon_barfwc", "popup_barfwc",
                                            i18n_barfwc + ": " + i18n_na);
                                    }
                                } else if (resp.accessibility.hasOwnProperty("barf_lift")) {
                                    let i18n_barflift = "Barrierefreier Aufzug";
                                    if (resp.accessibility.barf_lift === "FULLY_BARF") {
                                        cn.accessibilityIconYes("icon_barflift", "popup_barflift",
                                            i18n_barflift + ": " + i18n_yes);
                                    } else if (resp.accessibility.barf_lift === "SEMI_BARF") {
                                        cn.accessibilityIconPartly("icon_barflift", "popup_barflift",
                                            i18n_barflift + ": " + i18n_partly);
                                    } else if (resp.accessibility.barf_lift === "NOT_BARF") {
                                        cn.accessibilityIconNo("icon_barflift", "popup_barflift",
                                            i18n_barflift + ": " + i18n_no);
                                    } else if (resp.accessibility.barf_lift === "NO_INFO") {
                                        cn.accessibilityIconNa("icon_barflift", "popup_barflift",
                                            i18n_barflift + ": " + i18n_na);
                                    }
                                } else if (resp.accessibility.hasOwnProperty("barf_recreationroom")) {
                                    let i18n_recreationroom = "Ruheraum";
                                    if (resp.accessibility.barf_recreationroom === "FULLY_BARF") {
                                        cn.accessibilityIconYes("icon_recreationroom", "popup_recreationroom",
                                            i18n_recreationroom + ": " + i18n_yes);
                                    } else if (resp.accessibility.barf_recreationroom === "SEMI_BARF") {
                                        cn.accessibilityIconPartly("icon_recreationroom", "popup_recreationroom",
                                            i18n_recreationroom + ": " + i18n_partly);
                                    } else if (resp.accessibility.barf_recreationroom === "NOT_BARF") {
                                        cn.accessibilityIconNo("icon_recreationroom", "popup_recreationroom",
                                            i18n_recreationroom + ": " + i18n_no);
                                    } else if (resp.accessibility.barf_recreationroom === "NO_INFO") {
                                        cn.accessibilityIconNa("icon_recreationroom", "popup_recreationroom",
                                            i18n_recreationroom + ": " + i18n_na);
                                    }
                                } else if (resp.accessibility.hasOwnProperty("barf_lecturerarea")) {
                                    let i18n_wheelchairlecturer = "Barrierefreie Dozentenzone";
                                    if (resp.accessibility.barf_lecturerarea === "FULLY_BARF") {
                                        cn.accessibilityIconYes("icon_wheelchairlecturer", "popup_wheelchairlecturer",
                                            i18n_wheelchairlecturer + ": " + i18n_yes);
                                    } else if (resp.accessibility.barf_lecturerarea === "SEMI_BARF") {
                                        cn.accessibilityIconPartly("icon_wheelchairlecturer", "popup_wheelchairlecturer",
                                            i18n_wheelchairlecturer + ": " + i18n_partly);
                                    } else if (resp.accessibility.barf_lecturerarea === "NOT_BARF") {
                                        cn.accessibilityIconNo("icon_wheelchairlecturer", "popup_wheelchairlecturer",
                                            i18n_wheelchairlecturer + ": " + i18n_no);
                                    } else if (resp.accessibility.barf_lecturerarea === "NO_INFO") {
                                        cn.accessibilityIconNa("icon_wheelchairlecturer", "popup_wheelchairlecturer",
                                            i18n_wheelchairlecturer + ": " + i18n_na);
                                    }

                                    let i18n_wheelchairspace = "Rollstuhlplätze";
                                    if (resp.accessibility.barf_wheelchairspace === "FULLY_BARF") {
                                        cn.accessibilityIconYes("icon_wheelchairtable", "popup_wheelchairtable",
                                            i18n_wheelchairspace + ": " + +i18n_yes);
                                    } else if (resp.accessibility.barf_wheelchairspace === "SEMI_BARF") {
                                        cn.accessibilityIconPartly("icon_wheelchairtable", "popup_wheelchairtable",
                                            i18n_wheelchairspace + ": " + i18n_partly);
                                    } else if (resp.accessibility.barf_wheelchairspace === "NOT_BARF") {
                                        cn.accessibilityIconNo("icon_wheelchairtable", "popup_wheelchairtable",
                                            i18n_wheelchairspace + ": " + i18n_no);
                                    } else if (resp.accessibility.barf_wheelchairspace === "NO_INFO") {
                                        cn.accessibilityIconNa("icon_wheelchairtable", "popup_wheelchairtable",
                                            i18n_wheelchairspace + ": " + i18n_na);
                                    }

                                    let i18n_markedsteps = "Antritts-Endstufen markiert";
                                    if (resp.accessibility.barf_markedsteps === "FULLY_BARF") {
                                        cn.accessibilityIconYes("icon_markedsteps", "popup_markedsteps",
                                            i18n_markedsteps + ": " + i18n_yes);
                                    } else if (resp.accessibility.barf_markedsteps === "SEMI_BARF") {
                                        cn.accessibilityIconPartly("icon_markedsteps", "popup_markedsteps",
                                            i18n_markedsteps + ": " + i18n_partly);
                                    } else if (resp.accessibility.barf_markedsteps === "NOT_BARF") {
                                        cn.accessibilityIconNo("icon_markedsteps", "popup_markedsteps",
                                            i18n_markedsteps + ": " + i18n_no);
                                    } else if (resp.accessibility.barf_markedsteps === "NO_INFO") {
                                        cn.accessibilityIconNa("icon_markedsteps", "popup_markedsteps",
                                            i18n_markedsteps + ": " + i18n_na);
                                    }
                                }
                            } else {
                                noinfo = true;
                            }
                        } else {
                            noinfo = true;
                        }
                        
                        if (noinfo) {
                            cn.accessibilityIconNa("icon_wheelchair", "popup_wheelchair", "Barrierefreier Zugang" + ": " + i18n_na);
                            cn.accessibilityIconNa("icon_barfwc", "popup_barfwc", "Barrierefreies WC" + ": " + i18n_na);
                            cn.accessibilityIconNa("icon_barflift", "popup_barflift", "Aufzug" + ": " + i18n_na);
                            cn.accessibilityIconNa("icon_recreationroom", "popup_recreationroom", "Ruheraum" + ": " + i18n_na);
                            cn.accessibilityIconNa("icon_wheelchairlecturer", "popup_wheelchairlecturer", "Barrierefreie Dozentenzone" + ": " + i18n_na);
                            cn.accessibilityIconNa("icon_wheelchairtable", "popup_wheelchairtable", "Rollstuhlplätze" + ": " + i18n_na);
                            cn.accessibilityIconNa("icon_markedsteps", "popup_markedsteps", "Antritts-Endstufen markiert" + ": " + i18n_na);
                        }

                        if (resp.routing) {
                            cn.$("routing_buttons").style.display = "block";
                        }

                        if (resp.links && !cn.isEmpty(resp.links)) {
                            let linkList = cn.$("link_list");
                            for (let link in resp.links) {
                                let li = document.createElement("LI");
                                li.className = "reg";
                                let a = document.createElement("A");
                                a.href = link;
                                a.innerText = resp.links[link];
                                li.appendChild(a);
                                linkList.appendChild(li);
                            }
                        }
                    });
                    cn.AJAX.callApi("GET", "roominfo/542100.2130?all=true", "", torun);
                }

                function updateGlobalTilePos() {
                    g_tilex = ((data_canv_width * quali_steps[quali_cur]) / subpics_size) * 0.5;
                    g_tiley = ((data_canv_height * quali_steps[quali_cur]) / subpics_size) * 0.5;
                    g_xtrans = g_tilex - Math.floor(g_tilex);
                    g_ytrans = g_tiley - Math.floor(g_tiley);
                    g_tilex = Math.floor(g_tilex);
                    g_tiley = Math.floor(g_tiley);
                }


                function drawAllVisibleSubpics() {

                    var cq_pixel_size = (subpics_size / quali_steps[quali_cur]) * scale;

                    var xinview = new Array(2);
                    var yinview = new Array(2);

                    var tmp = g_tilex - Math.ceil((layx - g_xtrans * cq_pixel_size) / cq_pixel_size);
                    if (tmp < 0)
                        tmp = 0;
                    xinview[0] = tmp;

                    tmp = g_tilex + Math.ceil((canvas_width - layx - (1 - g_xtrans) * cq_pixel_size) / cq_pixel_size);
                    if (tmp > subpics[quali_cur].h - 1)
                        tmp = subpics[quali_cur].h - 1;
                    xinview[1] = tmp;

                    tmp = g_tiley - Math.ceil((layy - g_ytrans * cq_pixel_size) / cq_pixel_size);
                    if (tmp < 0)
                        tmp = 0;
                    yinview[0] = tmp;

                    tmp = g_tiley + Math.ceil((canvas_height - layy - (1 - g_ytrans) * cq_pixel_size) / cq_pixel_size);
                    if (tmp > subpics[quali_cur].v - 1)
                        tmp = subpics[quali_cur].v - 1;
                    yinview[1] = tmp;

                    visibleX = xinview;
                    visibleY = yinview;

                    // Bilder innerhalb des Sichtbereiches
                    for (var i = visibleX[0]; i <= visibleX[1]; i++) {
                        for (var j = visibleY[0]; j <= visibleY[1]; j++) {
                            var imgObj = new Image();
                            imgObj.xpos = i;
                            imgObj.ypos = j;
                            imgObj.onload = function () {
                                //Darstellen
                                addImageToCanvas(this, this.xpos, this.ypos);
                            }
                            imgObj.src = "images/etplan_cache/" + png_file_name + "_" + quali_steps[quali_cur] + "/" + i + "_" + j + ".png/nobase64";
                        }
                    }
                }

                function rescale() {
                    for (var i = 0; i < canvas.getChildren().length; i++) {
                        var laye = canvas.getChildren()[i];
                        laye.scale({x: scale, y: scale});
                    }
                }

                function redrawAllLayers() {
                    for (var i = 0; i < canvas.getChildren().length; i++) {
                        var laye = canvas.getChildren()[i];
                        laye.batchDraw();
                    }
                }

                function addImageToCanvas(imgobj, x, y) {
                    var curquali_size = subpics_size / quali_steps[quali_cur];
                    var img = new Kinetic.Image({
                        x: (x - g_tilex) * curquali_size - g_xtrans * curquali_size,
                        y: (y - g_tiley) * curquali_size - g_ytrans * curquali_size,
                        width: subpics_size,
                        height: subpics_size,
                        scale: {x: 1 / (quali_steps[quali_cur]), y: 1 / (quali_steps[quali_cur])},
                        opacity: 1,
                        listening: false,
                        visible: true,
                        image: imgobj
                    });
                    cache.add(img);
                    cache.batchDraw();
                }

                function move(deltaX, deltaY) {
                    layx += deltaX;
                    layy += deltaY;
                    for (var i = 0; i < canvas.getChildren().length; i++) {
                        var laye = canvas.getChildren()[i];
                        laye.move({x: deltaX, y: deltaY});
                    }
                }

                // Raumbelegungsplan zeugs ab hier
                function initMoreLSK() {
                    var w1 = cn.$("woche1");
                    var w2 = cn.$("woche2");
                    if (w1) {
                        w1 = cn.getElements(w1, "tbody-tr-td-span");
                        for (var i = 0; i < w1.length; i++) {
                            if (w1[i].className === "more") {
                                w1[i].parentNode.addEventListener("mousemove", RAplan.showMore, false);
                                w1[i].parentNode.addEventListener("mouseleave", RAplan.hideMore, false);
                            }
                        }
                    }
                    if (w2) {
                        w2 = cn.getElements(w2, "tbody-tr-td-span");
                        for (var i = 0; i < w2.length; i++) {
                            if (w2[i].className === "more") {
                                w2[i].parentNode.addEventListener("mousemove", RAplan.showMore, false);
                                w2[i].parentNode.addEventListener("mouseleave", RAplan.hideMore, false);
                            }
                        }
                    }
                }

                this.showMore = function (event) {
                    var x = event.clientX - cn.$("main").getBoundingClientRect().left + 20;
                    var y = event.clientY - cn.$("main").getBoundingClientRect().top + 20 + cn.$("main").scrollTop;
                    var mi = cn.$("mouseinfo");
                    var more = cn.getElements(event.target, "span");
                    for (var i = 0; i < more.length; i++) {
                        if (more[i].className === "more") {
                            mi.innerHTML = more[i].innerHTML;
                            break;
                        }
                    }
                    mi.style.left = x + "px";
                    mi.style.top = y + "px";
                    mi.style.display = "block";
                };

                this.hideMore = function () {
                    cn.$("mouseinfo").style.display = "none";
                };
                
                /**
                 * Set the value of checkbox if only body is to be printed.
                 * @returns {undefined}
                 */
                this.setBodyOnlyCheckbox = function () {
                    if (document.getElementById("body_only").checked) {
                        document.getElementById("body_only").value = "true";
                    } else {
                        document.getElementById("body_only").value = "false";
                    }     
                };
                
                /**
                 * Set the value of brass checkbox and based on choice generate different web formula
                 * in case of office doorplate.
                 * @returns {undefined}
                 */
                this.setBrassCheckbox = function () {
                    if (document.getElementById("is_brass").checked) {
                        document.getElementById("is_brass").value = "true";
                        document.getElementById("body_only").disabled = false;
                        document.getElementById("bodyOnly_label").style.color = "#000000"; //black
                    } else {
                        document.getElementById("is_brass").value = "false";
                        document.getElementById("body_only").checked = false ;
                        document.getElementById("body_only").disabled= true;
                        document.getElementById("bodyOnly_label").style.color = "#aaa";
                    }
                    let dpTemplate = cn.$("doorplate_template").selectedIndex;
                    let totalNameboxes = document.getElementsByClassName('name_box').length - 1 ;
                    let buttons = cn.$("add_remove_buttons");
                     
                   //Generating web formula for brass version of office doorplates
                    if (dpTemplate === 1 && document.getElementById("is_brass").checked && totalNameboxes > 1) {
                        if (totalNameboxes > 2) {
                            buttons.children[0].style.display = "none";
                            buttons.children[0].disabled = true;
                        } else {
                            buttons.children[0].style.display = "";
                            buttons.children[0].disabled = false;
                        } 
                        for (let i = 0; i <= totalNameboxes; i++) {
                            document.getElementById("textareaTitle_" + i).style.display = "none";
                            document.getElementById("textarea_" + i).value = "";
                            document.getElementById("textarea_" + i).style.display = "none";
                            document.getElementById("textarea_" + i).disabled = true;
                            if (totalNameboxes > 2) {
                                for (let j = 0; j <= totalNameboxes; j++) {
                                    document.getElementById("functionTitle_" + j).style.display = "none";
                                    document.getElementById("function_" + j).value = "";
                                    document.getElementById("function_" + j).style.display = "none";
                                    document.getElementById("function_" + j).disabled = true;
                                    if (totalNameboxes > 3) {        
                                        for (let k = 4; k <= totalNameboxes; k++) {
                                            var name_box = document.getElementById("name_box_" + k);
                                            document.getElementById("pre_title_" + k).value = "";
                                            document.getElementById("name_" + k).value = "";
                                            document.getElementById("post_title_" + k).value = "";
                                            name_box.style.display = "none";
                                            name_box.disabled = true;
                                       }
                                    } 
                               }
                            } 
                        }
                    }
                    
                    //switching back to normal office dp on uncheck
                    if (dpTemplate === 1 && document.getElementById("is_brass").checked === false && totalNameboxes > 1) {
                        if (totalNameboxes > 4) {
                            buttons.children[0].style.display = "none";
                            buttons.children[0].disabled = true;
                        } else {
                            buttons.children[0].style.display = "";
                            buttons.children[0].disabled = false;
                        }  
                         for (let i = 0; i <= totalNameboxes; i++) {
                            if (document.getElementById("textarea_" + i).disabled) {
                                document.getElementById("textareaTitle_" + i).style.display = "";
                                document.getElementById("textarea_" + i).style.display = "";
                                document.getElementById("textarea_" + i).disabled = false;
                            }
                            if (totalNameboxes > 2) {
                                for (let j = 0 ; j <= totalNameboxes; j++) {
                                    if (document.getElementById("function_" + j).disabled) {
                                        document.getElementById("functionTitle_" + j).style.display = "";
                                        document.getElementById("function_" + j).style.display = "";
                                        document.getElementById("function_" + j).disabled = false;
                                        document.getElementById("textareaTitle_" + j).style.display = "none";
                                        document.getElementById("textarea_" + j).value = "";
                                        document.getElementById("textarea_" + j).style.display = "none";
                                        document.getElementById("textarea_" + j).disabled = true;
                                    }
                                    if (totalNameboxes > 3) {
                                        document.getElementById("functionTitle_" + j).style.display = "none";
                                        document.getElementById("function_" + j).style.display = "none";
                                        document.getElementById("function_" + j).disabled = true;
                                        for (let k = 0 ; k <= totalNameboxes; k++) {
                                            if (document.getElementById("name_box_" + k).disabled) {
                                               document.getElementById("name_box_" + k).style.display = "";
                                               document.getElementById("name_box_" + k).disabled = false;
                                            }
                                        }
                                    }
                                }
                            }   
                        }
                    }
                };
                
                /**
                 * Put remaining characters value into a counter element.
                 * @param {Number} currentChars number of current character
                 * @param {Number} maxChars number of maximum characters
                 * @param {Number} idx the index of the counter which holds the number of remaining characters
                 */
                let characterCounter = function (currentChars, maxChars, idx) {
                    var remainChars = maxChars - currentChars;
                    if (remainChars < 0) {
                        document.getElementById("counter_" + idx).innerHTML = "("+0+"):";
                    } else {
                        document.getElementById("counter_" + idx).innerHTML = "("+remainChars+"):";
                    }
                };
                
                /**
                 * Retrieve the doorplate editor from the server and show it instead of
                 * the digital doorplate (legacy doorplates).
                 * Check access permissions on doorplate POST endpoint, then retrieve already filled in data and show
                 * editor (BLO doorplates).
                 * @param {boolean} blo whether this is a BLO doorplate or not
                 * @@param {boolean} templateChanged whether user has changed template or not  
                 * @returns {undefined}
                 */
                this.editDoorplate = function (blo, templateChanged) {
                    if (blo) {
                        let torun = new Map();
                        torun.set(204, function (responseText) {
                            let torun1 = new Map();
                            torun1.set(200, function (responseText) {
                                cn.$("doorplatecontent").style.display = "none";
                                let response = JSON.parse(responseText);
                                let resp = response.doorplate;
                                let isBrass = response.doorplate.brassTemplate;
                                cn.$("is_brass").checked = isBrass;
                                RAplan.setBrassCheckbox();let centralRoomIcons = resp.centralRoomIcons;
                                let functionRoomIcons = resp.functionRoomIcons;
                                let customOptionVal = "other";
                                if (templateChanged) {
                                    let plateType = cn.$("doorplate_template").selectedIndex;
                                    switchDoorplateTemplate(true, plateType, centralRoomIcons, functionRoomIcons,
                                        customOptionVal);
                                } else {
                                    let plateType = resp.template;
                                    cn.$("doorplate_template")[resp.template].selected = true;
                                    switchDoorplateTemplate(false, plateType, centralRoomIcons, functionRoomIcons,
                                        customOptionVal);
                                    if (resp.costcenterDe !== null && resp.costcenterDe.length > 0) {
                                        cn.$("costcenter_de").value = resp.costcenterDe;
                                    }
                                    if (resp.costcenterEn !== null && resp.costcenterEn.length > 0) {
                                        cn.$("costcenter_en").value = resp.costcenterEn;
                                    }

                                    let container = cn.$("doorplateeditor_generated");
                                    let buttons = cn.$("add_remove_buttons");

                                    if (plateType === 1 || plateType === 2) {
                                        for (let i = 0; i < resp.persons.length; i++) {
                                            if (!cn.$("name_" + i)) {
                                                addNameToDoorplate(container, i, buttons, isBrass);
                                            }
                                            cn.$("pre_title_" + i).value = resp.persons[i].preTitle;
                                            cn.$("post_title_" + i).value = resp.persons[i].postTitle;
                                            cn.$("name_" + i).value = resp.persons[i].name;
                                            if ((i < 6 && !isBrass) || (i < 3 && isBrass)) {
                                                cn.$("function_" + i).value = resp.persons[i].function;
                                            }
                                            if ((i < 3 && !isBrass) || (i < 2 && isBrass)) {
                                                cn.$("textarea_" + i).value = resp.persons[i].textarea;
                                            }
                                            cn.$("consent_" + i).checked = resp.persons[i].consent;
                                        }
                                    } else if (plateType === 3) {
                                        cn.$("headline").value = resp.headline;
                                        cn.$("textarea_func").value = resp.textarea;
                                    }
                                    if (plateType === 1) {
                                        buttons.children[0].disabled = false;
                                        buttons.children[0].style.display = "";
                                    } else {
                                        buttons.children[0].disabled = true;
                                        buttons.children[0].style.display = "none";
                                        buttons.children[1].disabled = true;
                                        buttons.children[1].style.display = "none";
                                    }

                                    let doors = response.doors;
                                    for (let i = 0; i < doors.length; i++) {
                                        let option = document.createElement("OPTION");
                                        option.value = doors[i].doorid;
                                        if (response.roomid === doors[i].room1) {
                                            option.innerHTML = doors[i].room2;
                                        } else {
                                            option.innerHTML = doors[i].room1;
                                        }
                                        cn.$("sel_door").appendChild(option);
                                    }
                                    if (doors.length < 2) {
                                        cn.$("sel_door").disabled = true;
                                    }
                                    cn.$("doorplateeditor_form").style.display = "block";
                                }
                                if (resp.roomFunction !== null) {
                                    cn.switchToOption(cn.$("room_function"), resp.roomFunction);
                                }
                                if (resp.customRoomFunctionDe !== null || resp.customRoomFunctionEn !== null) {
                                    cn.switchToOption(cn.$("room_function"), customOptionVal);
                                    RAplan.customRoomFunction(customOptionVal,
                                        resp.customRoomFunctionDe !== null ? resp.customRoomFunctionDe : "",
                                        resp.customRoomFunctionEn !== null ? resp.customRoomFunctionEn : "");
                                }
                            });
                            torun1.set(500, function (responseText) {
                                alert("An error occurred during PDF generation.");
                            });
                            cn.AJAX.callApi("GET", "roominfo/542100.2130?doorplate=true", "", torun1);
                        });
                        torun.set(400, function (responseText) {
                            alert("Invalid room ID.");
                        });
                        torun.set(403, function (responseText) {
                            cn.deleteCookie("loginToken");
                            cn.$("post_login_call").onclick = function () {
                                RAplan.editDoorplate(true);
                            };
                            cn.$("post_login_reload").checked = false;
                            cn.$("modalloginformdiv").style.display = "block";
                        });
                        cn.AJAX.callApi("HEAD", "roominfo/542100.2130/doorplate", "", torun);
                    } else {    // show non-BLO doorplate editor:
                        cn.AJAX.send("raum/542100.2130/editdoorplate", "blo=" + blo, function (responseText) {
                            cn.$("doorplatecontent").innerHTML = responseText;
                            if (!responseText.startsWith("<form")) {
                                cn.$("modalloginformdiv").style.display = "block";
                            }
                        });
                    }
                };

                /**
                 * Collects data from BLO doorplate editor form and sends it to savedoorplate-
                 * endpoint. If the response says there are errors, then display the error,
                 * otherwise download the doorplate PDF from the server.
                 * @returns {undefined} nothing
                 */
                this.saveBloDoorplate = function () {
                    let params = {};
                    try {
                        let form = cn.$("doorplateeditor_form").elements;
                        let personsArray = [];
                        for (let i = 0; i < 6; i++) {
                            if (cn.$("name_" + i)) {
                                let person = {
                                    consent: cn.$("consent_" + i).checked,
                                    preTitle: cn.$("pre_title_" + i).value,
                                    name: cn.$("name_" + i).value,
                                    postTitle: cn.$("post_title_" + i).value,
                                    function: cn.$("function_" + i) !== null ? cn.$("function_" + i).value : null,
                                    textarea: cn.$("textarea_" + i) !== null ? cn.$("textarea_" + i).value : null,
                                };
                                personsArray.push(person);
                            }
                        }

                        params = {
                            template: form.doorplate_template.value,
                            blo: true,
                            costcenterDe: form.costcenter_de.value,
                            costcenterEn: form.costcenter_en.value,
                            textarea: form.textarea_func !== undefined ? form.textarea_func.value : null,
                            roomFunction: form.room_function !== undefined ? form.room_function.value : null,
                            customRoomFunctionDe: form.custom_function_de !== undefined ? form.custom_function_de.value : null,
                            customRoomFunctionEn: form.custom_function_en !== undefined ? form.custom_function_en.value : null,
                            headline: form.headline !== undefined ? form.headline.value : null,
                            brassTemplate: form.is_brass.checked,
                            brassBodyOnly: form.body_only.checked,
                            door: form.sel_door.value !== undefined && form.sel_door.value !== "" ? form.sel_door.value : null,
                            persons: personsArray
                        };
                    } catch (ex) {
                        cn.reportError(ex, "saveBloDoorplate()");
                        return;
                    }
                    let torun = new Map();
                    torun.set(200, function (responseText) {
                        let resp = JSON.parse(responseText);
                        cn.downloadFile(resp.filename);
                    });
                    torun.set(400, function (responseText) {
                        alert("Can't create doorplate due to input problems.\n" + responseText.message);
                    });
                    torun.set(403, function (responseText) {
                        cn.deleteCookie("loginToken");
                        cn.$("post_login_call").onclick = function () {
                            RAplan.saveBloDoorplate();
                        };
                        cn.$("post_login_reload").checked = false;
                        cn.$("modalloginformdiv").style.display = "block";
                    });
                    torun.set(500, function (responseText) {
                        alert("Can't create doorplate due to server problems.\n" + responseText.message);
                    });
                    cn.AJAX.callApi("POST", "roominfo/542100.2130/doorplate", JSON.stringify(params), torun);
                };

                /**
                 * Counts existing name rows and sets the global index_names variable to the
                 * next index. Adds new rows for function and name afterwards. Previous
                 * function will be prefilled.
                 * @returns {undefined} nothing
                 */
                this.addFunctionNameRows = function () {
                    for (var i = 0; i < 10; i++) {
                        if (cn.$("name_" + i) !== null)
                            var idx = i;
                        else
                            break;
                    }
                    if (idx < 9) {
                        idx++;
                        var newfunctionrow = document.createElement("tr");
                        newfunctionrow.innerHTML = "<td>"
                                + "Funktion"
                                + ":</td><td><input id=\"function_" + idx + "\""
                                + " maxlength=\"64\" type=\"text\" name=\"function_" + idx + "\""
                                + " value=\"" + cn.$("function_" + (idx - 1)).value + "\""
                                + "></td>";
                        var newnamerow = document.createElement("tr");
                        newnamerow.innerHTML = "<td>"
                                + "Name"
                                + ":</td><td><input id=\"name_" + idx + "\""
                                + " maxlength=\"64\" type=\"text\" name=\"name_" + idx + "\""
                                + "></td><td><input style=\"width:1%;\" type=\"checkbox\" name=\"consent_" + idx + "\""
                                + ">*</td>";

                        var buttonrow = cn.$("row_add_name_button");
                        buttonrow.parentNode.insertBefore(newfunctionrow, buttonrow);
                        buttonrow.parentNode.insertBefore(newnamerow, buttonrow);
                    }
                };

                /**
                 * Loops through all name_box entries present on website to get the index of the last one.
                 * @return {Number} the index of the last "name_box_" entry
                 */
                 function lastNameBoxIdx() {
                    for (let i = 0; i < 7; i++) {
                        if (cn.$("name_box_" + i) === null) {
                            return i - 1;
                        }
                    }
                }
                
                /**
                 * Adds pre-title, name, post-title, function and text field input to a div with the ID 
                 * name_box_idx and adds it to the overhanded container.
                 * Does this only for 0 <= idx <= 6, so no more than six entries are possible.
                 * controls visibility of add and remove buttons for adding and removing name-input elements.
                 * @param {HTML-div} container gets the labels and inputs as children
                 * @param {Number} idx the index of the name box that will be created
                 * @param {HTML-div} button_div div to hold add and remove buttons
                 * @param {boolean} isBrass whether the name is added to a brass doorplate or not.
                 * This influences the number of additional entries per name like function and textarea.
                 * @return {undefined} nothing
                 */
                 function addNameToDoorplate(container, idx, button_div, isBrass) {
                    let maxChar = 128;
                    if (idx > 0) {
                        maxChar = 52;
                    } else {
                        maxChar = 128;
                        }
                    if ((0 <= idx && idx <= 6 && !isBrass) || (0 <= idx && idx <= 4 && isBrass)) {
                        let name_box = document.createElement("DIV");
                        name_box.id = "name_box_" + idx;
                        name_box.className = "name_box";

                        let div = document.createElement("DIV");
                        div.className = "dpedit_left";
                        div.innerHTML = "Akad. Titel:";
                        name_box.appendChild(div);
                        let indiv = document.createElement("INPUT");
                        indiv.className = "dpedit_mid";
                        indiv.id = "pre_title_" + idx;
                        indiv.type = "text";
                        indiv.maxLength = "128";
                        indiv.name = "pre_title_" + idx;
                        name_box.appendChild(indiv);

                        div = document.createElement("DIV");
                        div.className = "dpedit_left";
                        div.innerHTML = "Name:";
                        name_box.appendChild(div);
                        indiv = document.createElement("INPUT");
                        indiv.className = "dpedit_mid";
                        indiv.id = "name_" + idx;
                        indiv.type = "text";
                        indiv.maxLength = "128";
                        indiv.name = "name_" + idx;
                        name_box.appendChild(indiv);
                        
                        let inner_div = document.createElement("DIV");
                        inner_div.className = "dpedit_right";
                        let consent = document.createElement("INPUT");
                        consent.id = "consent_" + idx;
                        consent.type = "checkbox";
                        consent.name = "consents";
                        consent.style.width = "auto";
                        consent.style.marginLeft = "6px";
                        consent.content = "*";
                        inner_div.appendChild(consent);
                        let consent_label = document.createElement("LABEL");
                        consent_label.for = "consent_" + idx;
                        consent_label.innerHTML = "*";
                        inner_div.appendChild(consent_label);
                        name_box.appendChild(inner_div);
                        
                        div = document.createElement("DIV");
                        div.className = "dpedit_left";
                        div.innerHTML = "Akad. Titel (nachstehend):";
                        name_box.appendChild(div);
                        indiv = document.createElement("INPUT");
                        indiv.className = "dpedit_mid";
                        indiv.id = "post_title_" + idx;
                        indiv.type = "text";
                        indiv.maxLength = "128";
                        indiv.name = "post_title_" + idx;
                        name_box.appendChild(indiv);

                        if ((idx < 6 && !isBrass) || (idx < 3 && isBrass)) {      // add functions if less than six
                            div = document.createElement("DIV");
                            div.id = "functionTitle_" + idx;
                            div.className = "dpedit_left";
                            div.innerHTML = "Funktion:";
                            name_box.appendChild(div);
                            indiv = document.createElement("INPUT");
                            indiv.className = "dpedit_mid";
                            indiv.id = "function_" + idx;
                            indiv.type = "text";
                            indiv.maxLength = "128";
                            indiv.name = "function_" + idx;
                            name_box.appendChild(indiv);
                        }

                        if ((idx < 3 && !isBrass) || (idx < 2 && isBrass)) {      // add textarea if less than four entries
                            div = document.createElement("DIV");
                            div.id = "textareaTitle_" + idx;
                            div.className = "dpedit_left";
                            div.innerHTML = "Freitextfeld";
                            let counter_span = document.createElement("SPAN");
                            counter_span.id = "counter_" + idx;
                            counter_span.innerHTML = " ("+maxChar+"):";
                            div.appendChild(counter_span);
                            name_box.appendChild(div);
                            indiv = document.createElement("TEXTAREA");
                            indiv.className = "dpedit_mid";
                            indiv.id = "textarea_" + idx;
                            indiv.type = "text";
                            indiv.maxLength = maxChar;
                            indiv.name = "textarea_" + idx;
                            indiv.onkeyup = function() {
                                characterCounter(indiv.value.length, maxChar, idx);
                            };
                            name_box.appendChild(indiv);     
                        }

                        container.appendChild(name_box);
                    }
                    if ((idx <= 4 && !isBrass) || (idx < 3 && isBrass)) {
                        button_div.children[0].style.display = "";
                        button_div.children[0].disabled = false;
                        button_div.children[1].style.display = "";
                        button_div.children[1].disabled = false;
                    } else {
                        button_div.children[0].style.display = "none";
                        button_div.children[0].disabled = true;
                    }
                    
                    if ((idx > 3 && !isBrass) || (idx > 2 && isBrass)) {
                        for (let i = 0; i <= idx; i++) {
                            if (cn.$("functionTitle_" + i)) {
                                cn.$("functionTitle_"+i).style.display = "none";
                            }
                            if (cn.$("function_" + i)) {
                                cn.$("function_"+i).value = "";
                                cn.$("function_"+i).style.display = "none";
                                cn.$("function_"+i).disabled = true;
                            }
                        }
                    }
                    if ((idx > 2 && !isBrass) || (idx > 1 && isBrass)) {
                        for (let i = 0; i <= idx; i++) {
                            if (cn.$("textareaTitle_" + i)) {
                                cn.$("textareaTitle_"+i).style.display = "none";
                            }
                            if (cn.$("textarea_" + i)) {
                                cn.$("textarea_"+i).value = "";
                                cn.$("textarea_"+i).style.display = "none";
                                cn.$("textarea_"+i).disabled = true;
                              //  cn.$("counter_"+i).value = "";
                                cn.$("counter_"+i).style.display = "none";
                             //   cn.$("counter_"+i).disabled = true;
                            }
                        }
                    }
                }
                
                /**
                 * Removes the ith title-name-function-fastfact entry from the container as
                 * well as probably adding a textarea through addTextarea() and controlling
                 * the visibility of the add and remove buttons.
                 * @param {HTML-div} container the div to remove from
                 * @param {Number} idx
                 * @param {HTML-div} buttonDiv the div holding the add and remove button
                 * @param {boolean} brassOffice
                 * @returns {undefined} nothing
                 */
                 function removeNameFromDoorplate(container, idx, buttonDiv, brassOffice) {
                    if (cn.$("name_box_" + idx)) {
                        container.removeChild(cn.$("name_box_" + idx));
                    }
                    buttonDiv.children[0].style.display = "";
                    buttonDiv.children[0].disabled = false;

                    if (idx > 1) {
                        buttonDiv.children[1].style.display = "";
                        buttonDiv.children[1].disabled = false;
                    } else {
                        buttonDiv.children[1].style.display = "none";
                        buttonDiv.children[1].disabled = true;
                    } 
                    
                    if ((idx < 5 && brassOffice === false) || (idx < 4 && brassOffice === true)) { 
                        for (let i = idx - 1; i >= 0; i--) {
                            let div = document.getElementById("function_" + i);  
                            if (div.disabled) {
                                cn.$("functionTitle_" + i).style.display = "";
                                div.disabled = false;
                                div.style.display = "";
                            }
                            if ((idx < 4 && brassOffice === false) || (idx < 3 && brassOffice === true)) {
                                let div = document.getElementById("textarea_" + i);  
                                if (div.disabled) {
                                    cn.$("textareaTitle_" + i).style.display = "";
                                    div.disabled = false;
                                    div.style.display = "";
                                    cn.$("counter_" + i).style.display = "";
                                }
                            }
                       }
                    }
                }

                /**
                 * This function creates a skel for the doorplate type one has switched to. No values will be entered
                 * here. Cleans probably existing input fields and generates fresh ones depending on the chosen template.
                 * @param {boolean} templateChanged detects if template has been switched or not, could be needed in
                 * future for caching between doorplates switches
                 * @param {Number} plateType
                 * @param {Map} centralRoomIcons
                 * @param {Map} functionRoomIcons
                 * @param {String} customOptionVal the value for the custom/other option of room function select box
                 * @returns {undefined}
                 */
                 function switchDoorplateTemplate(templateChanged, plateType, centralRoomIcons, functionRoomIcons,
                                                  customOptionVal) {
                    let form = cn.$("doorplateeditor");
                    form.removeChild(cn.$("doorplateeditor_generated"));
                    if (cn.$("add_remove_buttons")) {
                        form.removeChild(cn.$("add_remove_buttons"));
                    }
                    
                    let changingContent = document.createElement("DIV");
                    changingContent.id = "doorplateeditor_generated";
                    form.insertBefore(changingContent, form.childNodes[form.childNodes.length - 2]);

                    let buttonDiv = document.createElement("DIV");
                    buttonDiv.id = "add_remove_buttons";
                    buttonDiv.style = "grid-column:2/2; justify-self:center;";

                    let add_button = document.createElement("INPUT");
                    add_button.type = "button";
                    add_button.value = "+";
                    add_button.style = "cursor:pointer; width:24px; height:100%; margin-right:2px;";

                    let remove_button = document.createElement("INPUT");
                    remove_button.type = "button";
                    remove_button.value = "-";
                    remove_button.style = "cursor:pointer; width:24px; height:100%; margin-left:2px;";

                    buttonDiv.appendChild(add_button);
                    buttonDiv.appendChild(remove_button);
                    form.insertBefore(buttonDiv, form.childNodes[form.childNodes.length - 2]);
                    
                    let customOption = document.createElement("OPTION");
                    customOption.value = customOptionVal;
                    customOption.innerHTML = "Andere";
                    
                    let Ddiv = document.getElementById("costcenter_de");
                    let Ediv = document.getElementById("costcenter_en");
                    let functionSelect = document.getElementById("room_function");
                                        
                    if (plateType === 1 || plateType === 2 ) {
                        cn.$("costcenter_de").required = true;
                        cn.$("costcenter_en").required = true;
                        addNameToDoorplate(changingContent, 0, buttonDiv, false);
                        cn.$("custom_functions").style.display = "none";
                        cn.$("custom_function_de").innerHTML = "";
                        cn.$("custom_function_en").innerHTML = "";
                    }
                    
                    if (plateType === 2) {                           //adding rome function to central and function DPs
                        if (functionSelect.childNodes.length > 1) {
                            functionSelect.innerHTML = '';
                        }
                        for (let key in centralRoomIcons){
                            let option = document.createElement("OPTION");
                            option.value = key ;
                            option.innerHTML = centralRoomIcons[key];
                            functionSelect.appendChild(option);
                        }                         
                        cn.sortSelect(functionSelect);
                        functionSelect.appendChild(customOption);
                        
                        add_button.disabled = true;
                        add_button.style.display = "none";
                        remove_button.disabled = true;
                        remove_button.style.display = "none"; 
                    }
                    
                    if (plateType === 3) {                                                    //function room goals
                        cn.$("costcenter_de").required = false;
                        cn.$("costcenter_en").required = false;
                        if (functionSelect.childNodes.length > 1) {
                            functionSelect.innerHTML = "";
                        }
                        cn.$("custom_functions").style.display = "none";
                        cn.$("custom_function_de").innerHTML = "";
                        cn.$("custom_function_en").innerHTML = "";

                        for (let key in functionRoomIcons){
                            let option = document.createElement("OPTION");
                            option.value = key ;
                            option.innerHTML = functionRoomIcons[key];
                            functionSelect.appendChild(option);
                        } 
                        cn.sortSelect(functionSelect);
                        functionSelect.appendChild(customOption);
                        
                        let box = document.createElement("DIV");
                        box.id = "headline_box";
                        box.className = "function_box";
                        let div = document.createElement("DIV");
                        div.className = "dpedit_left";
                        div.innerHTML = "Zusatzinfo:";
                        box.appendChild(div);
                        let indiv = document.createElement("INPUT");
                        indiv.className = "dpedit_mid";
                        indiv.id = "headline";
                        indiv.type = "text";
                        indiv.maxLength = "128";
                        indiv.name = "headline";
                        box.appendChild(indiv);
                        changingContent.appendChild(box);

                        box = document.createElement("DIV");
                        box.id = "textarea_box";
                        box.className = "function_box";
                        div = document.createElement("DIV");
                        div.className = "dpedit_left";
                        div.innerHTML = "Freitextfeld";
                        let counter_span = document.createElement("SPAN");
                        counter_span.id = "counter_0";
                        counter_span.innerHTML = " (256):";
                        div.appendChild(counter_span);
                        box.appendChild(div);
                        indiv = document.createElement("TEXTAREA");
                        indiv.className = "dpedit_mid";
                        indiv.id = "textarea_func";
                        indiv.type = "text";
                        indiv.maxLength = "256";
                        indiv.name = "textarea_func";
                        indiv.onkeyup = function() {
                            characterCounter(indiv.value.length, 256, 0);
                        };
                        box.appendChild(indiv);
                        changingContent.appendChild(box);
                        
                        cn.$("costcentreDe").style.display = "none";
                        cn.$("costcentreEn").style.display = "none";
                        Ddiv.style.display = "none" ;
                        Ediv.style.display = "none";
                        
                        add_button.disabled = true;
                        add_button.style.display = "none";
                        remove_button.disabled = true;
                        remove_button.style.display = "none";
                    } else {
                        cn.$("costcenter_de").required = true;
                        cn.$("costcenter_en").required = true;
                        cn.$("costcentreDe").style.display = "block";
                        cn.$("costcentreEn").style.display = "block";
                        Ediv.style.display = "block";
                        Ddiv.style.display = "block" ;
                    }
                    
                    if (plateType === 1) {              // removing room function from office rooms
                        cn.$("roomfunction").style.display = "none" ;
                        functionSelect.style.display = "none";
                    } else {
                        cn.$("roomfunction").style.display = "block" ;
                        functionSelect.style.display = "block";
                    }
                    
                    add_button.onclick = function () {
                        let namesLastIdx = lastNameBoxIdx();
                        let brassOffice = false ;
                        if (plateType === 1 && document.getElementById("is_brass").checked) {
                            brassOffice = true;
                        }
                        addNameToDoorplate(changingContent, ++namesLastIdx, buttonDiv, brassOffice);
                    };
                    remove_button.onclick = function () {
                        let names_last_idx = lastNameBoxIdx();
                        let brassOffice = false ;
                        if (plateType === 1 && document.getElementById("is_brass").checked) {
                            brassOffice = true;
                        }
                        removeNameFromDoorplate(changingContent, names_last_idx, buttonDiv, brassOffice);
                    };
                }

               /**
                * Generates input field for custom room function 
                * @param {String} optionVal the room function selected
                * @param {String} customDe the custom value in German to set, default is ""
                * @param {String} customEn the custom value in English to set, default is ""
                * @returns {undefined}
                */
                this.customRoomFunction = function (optionVal, customDe = "", customEn = "") {
                   if (optionVal == "other") {
                       cn.$("custom_function_de").value = customDe;
                       cn.$("custom_function_en").value = customEn;
                       cn.$("custom_functions").style.display = "grid";
                   } else {
                       cn.$("custom_function_de").value = "";
                       cn.$("custom_function_en").value = "";
                       cn.$("custom_functions").style.display = "none";
                   }
                };

                /**
                 * Construct and show edit dialog for selected data.
                 * @param {HTMLElement} div a HTML-DIV which containing the current value. The editbox will be shown below this
                 * div.
                 * @param {string} key the values key
                 * @param {string} valuetype the type of the edited data to use the correct input type. Can be "string",
                 * "number", "boolean", "date".
                 */
                let editBarfData = function (div, key, valuetype) {
                    if (lasteditbarf) {
                        lasteditbarf.style.backgroundColor = "#fff";
                    }
                    lasteditbarf = div;
                    
                    let editbox = cn.$("editbox");
                    div.parentNode.insertBefore(editbox, div.nextSibling);      // place editbox directly after div
                    div.style.backgroundColor = "#ddd";
                    let value_corrected = cn.$("value_corrected");
                    if (valuetype == "string") {
                        value_corrected.type = "text";
                    } else if (valuetype == "number") {
                        value_corrected.type = "number";
                    } else if (valuetype == "boolean") {
                        value_corrected.type = "checkbox";
                    } else if (valuetype == "date") {
                        value_corrected.type = "date";
                    }
                    let button = cn.$("send_corrected_value");
                    button.onclick = function () {
                        let value = cn.$("value_corrected").value;
                        if (valuetype == "number") {
                            value = parseInt(value, 10);
                        } else if (valuetype == "boolean") {
                            value = cn.$("value_corrected").checked;
                        }
                        let params = {
                            features: [
                                {
                                    key: key,
                                    value: value
                                }
                            ],
                            email: cn.$("feedbackmail").value
                        };
                        let torun = new Map();
                        torun.set(500, function (responseText) {
                            alert(JSON.parse(responseText).message);
                        });
                        torun.set(403, function (responseText) {
                            cn.deleteCookie("loginToken");
                            cn.$("post_login_reload").checked = false;  // dont reload page after successful login
                            cn.$("modalloginformdiv").style.display = "block";
                        });
                        torun.set(200, function (responseText) {
                            // close dialog on success
                            editbox.style.display = "none";
                            lasteditbarf.style.backgroundColor = "#fff";
                        });
                        cn.AJAX.callApi("POST", "roominfo/542100.2130", JSON.stringify(params), torun);
                    };
                    let closebutton = cn.$("closebutton");
                    closebutton.onclick = function () {
                        editbox.style.display = "none";
                        lasteditbarf.style.backgroundColor = "#fff";
                    };
                    
                    editbox.style.display = "inline";
                };
            };

            var Bar = new function () {

                var lines;

                this.init = function () {
                    lines = cn.getElements(cn.$("barf_info"), "tbody-tr");

                    for (var i = 0; i < lines.length; i++) {
                        var line = lines[i];
                        var divs = cn.getElements(cn.getElements(line, "td")[1], "div");
                        if (divs.length > 1) {
                            var klapp = cn.getElements(line, "td")[0];
                            var span = cn.getElements(klapp, "span")[0];
                            klapp.style.cursor = "pointer";
                            span.style.display = "block";
                            klapp.setAttribute("open", 1);
                            klapp.onclick = function () {
                                var klapp = this;
                                var span = cn.getElements(klapp, "span")[0];
                                var divs = cn.getElements(cn.getElements(klapp.parentNode, "td")[1], "div");
                                var open = parseInt(klapp.getAttribute("open"));
                                if (open === 1) {
                                    klapp.setAttribute("open", 0);
                                    span.style.backgroundPosition = "-290px 2px";
                                    for (var j = 1; j < divs.length; j++) {
                                        var div = divs[j];
                                        div.style.display = "block";
                                    }
                                } else {
                                    klapp.setAttribute("open", 1);
                                    span.style.backgroundPosition = "-267px 2px";
                                    for (var j = 1; j < divs.length; j++) {
                                        var div = divs[j];
                                        div.style.display = "none";
                                    }
                                }
                            };
                            for (var j = 1; j < divs.length; j++) {
                                var div = divs[j];
                                div.style.display = "none";
                            }
                        }
                    }
                };
            }

        </script>
    </head>
    <body>
        <!--
from https://www.w3schools.com/howto/tryit.asp?filename=tryhow_css_login_form_modal
Include this (right below <body>) via:

The following button can be used to activate this loginform:
<button onclick="cn.$('modalloginformdiv').style.display='block';" style="width:auto;">Login</button>
A display value for the university selector must be passed to modalloginform.jsp
via the jsp:param tag as shown above.
TUD is checked by default.
-->



<style>
    /* Full-width input fields */
    .logincontainer input[type=text], input[type=password] {
        width: 100%;
        padding: 12px 20px;
        margin: 8px 0;
        display: inline-block;
        border: 1px solid #ccc;
        box-sizing: border-box;
    }
    .logincontainer {
        padding: 16px;
    }
    /* The Modal (background) */
    .modal {
        display: none; /* Hidden by default */
        position: fixed; /* Stay in place */
        z-index: 9999; /* Sit on top */
        left: 0;
        top: 0;
        width: 100%; /* Full width */
        height: 100%; /* Full height */
        overflow: auto; /* Enable scroll if needed */
        background-color: rgb(0,0,0); /* Fallback color */
        background-color: rgba(0,0,0,0.4); /* Black w/ opacity */
        padding-top: 60px;
    }
    /* Modal Content/Box */
    .modal-content {
        background-color: #fefefe;
        margin: 5% auto 15% auto; /* 5% from the top, 15% from the bottom and centered */
        border: 1px solid #888;
        width: 80%; /* Could be more or less, depending on screen size */
        max-width: 412px;
    }
    /* Add Zoom Animation */
    .animate {
        -webkit-animation: animatezoom 0.6s;
        animation: animatezoom 0.6s
    }
    @-webkit-keyframes animatezoom {
        from {-webkit-transform: scale(0)}
        to {-webkit-transform: scale(1)}
    }
    @keyframes animatezoom {
        from {transform: scale(0)}
        to {transform: scale(1)}
    }
</style>
<div id="modalloginformdiv" class="modal">
    <form id="modalloginform" class="modal-content animate" action="" method="post" onsubmit="cn.Main.login(this); return false;">
        <div class="closecontainer">
            <button class="close" aria-label="Close loginform" type="button" onclick="cn.$('modalloginformdiv').style.display='none'">
                <span aria-hidden="true">&times;</span>
            </button>
        </div>

        <div class="logincontainer">
            <label id="input_user"><b>Benutzername</b></label>
            </br>
            <input type="text" placeholder="Enter Username" name="user" required aria-labelledby="input_user">
            </br>
            <label id="input_passwd"><b>Passwort</b></label>
            </br>
            <input type="password" placeholder="Enter Password" name="passwd" required aria-labelledby="input_passwd">
            </br>
            <div id="login_radio" style="display: none;">
                <input id="radio_tud" name="university" checked="" style="width: auto; background-color: white;" type="radio">
                <label for="radio_tud">TUD</label>
                <br>
                <input id="radio_htw" name="university" style="width: auto; background-color: white;" type="radio">
                <label for="radio_htw">HTW</label>
            </div>
            <input id="post_login_reload" style="display:none;" checked="checked" type="checkbox">  <!-- offers the possibility to control the post-login page reload, default is yes -->
            <input id="post_login_call" style="display:none;">  <!-- can store post-login-callback functions -->
            </br>
            <label id="login_error" style="display: none; color: red;"></label>

            <button id="loginbutton" class="login" type="submit">Login</button>
        </div>
    </form>
</div>
        <div id="portal-wrapper">
            




<div id="portal-top" >

    <a href="" title="Startseite des Campus Navigator TU Dresden" id="navigator-logo" aria-flowto="mainh1">
        Startseite des Campus Navigator TU Dresden
    </a>

    <div id="portal-globalnav">
        <div  id="tab_buildings">
            <a title="Karten, Gebäudeinformationen, Etagenpläne, Rauminformationen" href="">
                Gebäude
            </a>
        </div>
        <div  id="tab_facilities">
            <a title="Bereiche, Fakultäten und Einrichtungen" href="einrichtungen">
                Einrichtungen
            </a>
        </div>
<!--    <div  id="tab_routing">
            <a title="Routenplaner" href="routing">
                Routing
            </a>
        </div> -->
        <div > 
            <a title="Suche nach freien Räumen an TUD und HTW. Ergebnis des Projekts "Hochschulübergreifendes Flächenmanagement" (HÜFM)." href="huefm/start">
                Freiraumsuche
            </a>
        </div>
<!--         <div > 
            <a title="Campus Navigator - Mobile Apps" href="mobile_apps">
                Mobile Apps
            </a>
        </div> -->
        <div > 
            <a title="Hilfe zur Nutzung des Campus Navigator Dresden" href="hilfe">
                Hilfe
            </a>
        </div>
        <div > 
            <a title="Weitere Informationen über den Campus Navigator Dresden" href="ueber">
                Über
            </a>
        </div>
    </div>

    <div id="language_div">
        <div id='sprache'>Sprache</div>
        <div id='sprache_box'>
            <img src='images/flag_de.png' alt='' /><a href='' language='de'>Deutsch</a><br />
            <img src='images/flag_en.png' alt='' /><a href='' language='en'>Englisch</a><br />
            <img src='images/flag_pl.png' alt='' /><a href='' language='pl'>Polnisch</a><br />
            <img src='images/flag_cz.png' alt='' /><a href='' language='cz'>Tschechisch</a><br /> 
            <img src='images/flag_fr.png' alt='' /><a href='' language='fr'>Französisch</a><br />
            <img src='images/flag_es.png' alt='' /><a href='' language='es'>Spanisch</a><br />
            <img src='images/flag_ru.png' alt='' /><a href='' language='ru'>Russisch</a><br />
            <img src='images/flag_cn.png' alt='' /><a href='' language='cn'>Chinesisch</a><br />
        </div>
    </div>
</div>
<div id="portal-breadcrumbs" >
    <a href="" id='home_house'>
        <img src='images/home.png' alt="Startseite des Campus Navigator TU Dresden"/>
    </a>
    <span class='arrow'></span><a href=''>Gebäude</a><span class='arrow'></span><a href='karten/dresden'>Dresden Campus</a><span class='arrow'></span><a href='gebaeude/apb'>Andreas-Pfitzmann-Bau</a><span class='arrow'></span><a href='etplan/apb/00'>Etage 0</a><span class='arrow'></span><a href='raum/542100.2130' class='cur'>Raum E008</a>
    <div id="searchdiv">
        <a href="erweitertesuche" id="search_more" title="Erweiterte Suche" onclick='cn.SearchBox.fullSearch();return false;'><span class="hid">Erweiterte Suche</span></a>
        <div style="z-index:10001">
            <input type="text" id='search_assist' value='' aria-labelledby="searchboxtitle"/>
            <div id="searchborder">
                <input type="text" value='' id="searchbox" aria-labelledby="searchboxtitle"/>
                <img src="images/loading.gif" id="sb_loading_animation" alt="Ladeanimation" style="width:16px;"/>
                <div id="search_results"></div>
            </div>
        </div>
        <span id="searchboxtitle">Gebäude- und Raumsuche</span>
    </div>
</div>
                <div id="portal-main">
                    <div id="portal-menu_left">
                        <div id="menu_cont">
                            <h5>Andreas-Pfitzmann-Bau</h5>
                        <ul>
                            <li class='closed'><a href='' onclick='cn.Menu.menu_toggle_point(this);return false;'>Gebäude</a><ul class='submenu1'><li><a href='gebaeude/apb'>Übersicht</a></li><li><a href='barrierefrei/apb'>Barrierefreier Zugang</a></li></ul></li><li class='closed'><a href='' onclick='cn.Menu.menu_toggle_point(this);return false;'>Etagenpläne</a><ul class='submenu1'><li><a href='etplan/apb/03'>Etage 03</a></li><li><a href='etplan/apb/02'>Etage 02</a></li><li><a href='etplan/apb/01'>Etage 01</a></li><li><a href='etplan/apb/00'>Etage 00</a></li><li><a href='etplan/apb/-1'>Etage -1</a></li></ul></li><li><a href='hoersaele/apb'>Lehrräume</a></li>
                        </ul>
                    </div>
                </div>
                <div id="portal-menu_right">
                    <div id="menu_cont_right">
                        <h5 id="pos_in_building">Lage im Gebäude</h5>
                        <a href="etplan/apb/00/raum/542100.2130" aria-labelledby="pos_in_building"><div id="lageplan" style="border:1px solid #cccccc; margin:10px; width:200px; height:200px;"></div></a>
                        <h5>Informationen</h5>
                        <div style="margin:10px; width:200px;">
                            <h6>Raum E008</h6><p>Teilgebäude: </p><div class='p'>APB Andreas-Pfitzmann-Bau, Nöthnitzer Str. 46</div><p>Gebäude-Nr.: </p><div class='p'>5421</div><p>Etage: </p><div class='p'>Etage 0</div><p>Plätze: </p><div class='p'>30</div><p>Nutzung: </p><div class='p'>Übungs-/ Seminarraum mit DV</div><div style='height:9px; clear:both;'></div>
                            <div align="center" style="padding-right:30px;">
                                <div class="popup" onclick="cn.togglePopup('popup_wheelchair');">
                                    <object id="icon_wheelchair" type="image/svg+xml" data="images/symbols/icon_wheelchair.svg"></object>
                                    <span class="popuptext" id="popup_wheelchair"></span>
                                </div>
                                <div class="popup" onclick="cn.togglePopup('popup_barfwc');">
                                    <object id="icon_barfwc" type="image/svg+xml" data="images/symbols/icon_barfwc.svg"></object>
                                    <span class="popuptext" id="popup_barfwc"></span>
                                </div>
                                <br>
                                <div class="popup" onclick="cn.togglePopup('popup_barflift');">
                                    <object id="icon_barflift" type="image/svg+xml" data="images/symbols/icon_barflift.svg"></object>
                                    <span class="popuptext" id="popup_barflift"></span>
                                </div>
                                <div class="popup" onclick="cn.togglePopup('popup_wheelchairtable');">
                                    <object id="icon_wheelchairtable" type="image/svg+xml" data="images/symbols/icon_wheelchairtable.svg"></object>
                                    <span class="popuptext" id="popup_wheelchairtable"></span>
                                </div>
                                <div class="popup" onclick="cn.togglePopup('popup_wheelchairlecturer');">
                                    <object id="icon_wheelchairlecturer" type="image/svg+xml" data="images/symbols/icon_wheelchairlecturer.svg"></object>
                                    <span class="popuptext" id="popup_wheelchairlecturer"></span>
                                </div>
                                <br>
                                <div class="popup" onclick="cn.togglePopup('popup_markedsteps');">
                                    <object id="icon_markedsteps" type="image/svg+xml" data="images/symbols/icon_markedsteps.svg"></object>
                                    <span class="popuptext" id="popup_markedsteps"></span>
                                </div>
                                <div class="popup" onclick="cn.togglePopup('popup_recreationroom');">
                                    <object id="icon_recreationroom" type="image/svg+xml" data="images/symbols/icon_recreationroom.svg"></object>
                                    <span class="popuptext" id="popup_recreationroom"></span>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div id="main">
                    <div style="padding:14px;">
                        <h1 id="mainh1">APB(5421) E008 (APB/E008/U)</h1>
                        <!-- enable for other buildings and floors: -->
                        
                        <div>
                            <h2 style="float:left;">Digitales Türschild</h2>
                            <input class="editinputbutton" style="vertical-align:-23px;" value="edit" 
                                   title="Editieren" type="button" 
                                   onclick="RAplan.editDoorplate(true , false);" />
                            <br style="clear:both;" />
                        </div>
                        <div id="doorplatecontent">
                            <table><tbody><tr><td>Org.einheit* (deutsch):</td><td>Zentrale Räume Fakultät Informatik</td></tr><tr><td>Org.einheit* (englisch):</td><td></td></tr></tbody></table>
                        </div>
                        <form id="doorplateeditor_form" method="POST" onsubmit="RAplan.saveBloDoorplate();
                                return false;" style="display:none;">
                            <div id="doorplateeditor">
                                <p style="grid-column-start:1; grid-column-end:3;">Namen zu denen kein Einverständnis gegeben wurde tauchen nur im PDF auf und werden nicht gespeichert. Sie können Zeilenumbrüche für die Eingabefelder im PDF erzwingen, indem Sie an der gewünschten Position "\n" eingeben. Zeilenumbrüche innerhalb eines Wortes sind mit "-\n" möglich.</p>
                                <div class="dpedit_left">Raum:</div>
                                <div class="dpedit_mid">E008</div>
                                <div class="dpedit_left">Typ:</div>
                                <label for="is_brass" class="dpedit_mid">
                                    <input class="checkboxwithlabel" type="checkbox" id="is_brass" name="is_brass" onchange="RAplan.setBrassCheckbox()"/>
                                    Messingschild
                                    <input class="checkboxwithlabel" type="checkbox" id="body_only" name="body_only" disabled="disabled" style="margin-left:50px;" onchange="RAplan.setBodyOnlyCheckbox()"/>
                                    <label for="body_only" id="bodyOnly_label" style="color:#aaa";>Nur Unterteil</label>
                                </label>
                                <div class="dpedit_left">Vorlage:</div>
                                <select id="doorplate_template" class="dpedit_mid" name="doorplate_template" required onchange="RAplan.editDoorplate(true, true);">
                                    <option value="" selected="selected" disabled="disabled">Bitte Vorlage wählen</option> 
                                    <option value="1">Büroraum</option>
                                    <option value="2">Zentraler Raum der Organisationseinheit</option>
                                    <option value="3">Funktionsraum</option>
                                </select>
                                <div class="dpedit_left">Tür:</div>
                                <select id="sel_door" name="sel_door" class="dpedit_mid"></select>
                                <div class="dpedit_left" id="costcentreDe">Org.einheit* (deutsch):</div>
                                <input id="costcenter_de" class="dpedit_mid" name="costcenter_de" maxlength="128" required value="">
                                <div class="dpedit_left" id="costcentreEn">Org.einheit* (englisch):</div>
                                <input id="costcenter_en" class="dpedit_mid" name="costcenter_en" maxlength="128" required value="">
                                <div class="dpedit_left" id="roomfunction">Raumfunktion:</div>
                                <select id="room_function" name="room_function" class="dpedit_mid" onchange="RAplan.customRoomFunction(this.value)">
                                   <option value="" selected="selected" disabled="disabled">Bitte Raumfunktion wählen</option>  
                                </select>
                                <div id="custom_functions" style="grid-column-start:1;grid-column-end:3;display:none;grid-template-columns:200px 350px;grid-row-gap:5px;">
                                    <label for="custom_function_de" class="dpedit_left">Benutzerdefinierte Funktion (deutsch):</label>
                                    <input type="text" name="custom_function_de" id="custom_function_de" class="dpedit_mid"/>
                                    <label for="custom_function_en" class="dpedit_left">Benutzerdefinierte Funktion (englisch):</label>
                                    <input type="text" name="custom_function_en" id="custom_function_en" class="dpedit_mid"/>
                                </div>
                                <div id="doorplateeditor_generated"></div>

                                <input class="dpedit_left" type="submit" style="cursor:pointer; width:120%; height:20px;" value="Speichern & Türschild generieren">
                            </div>
                            <p>*) Hiermit willige ich freiwillig in die weltweite Veröffentlichung meines angegebenen Titels und vollen Namens sowie meiner Funktionsbezeichung auf der Raumseite des Campus Navigators ein. Mir ist bekannt, dass ich die Einwilligung ohne Rechtsfolgen verweigern oder mit Wirkung für die Zukunft jederzeit ohne Angabe von Gründen widerrufen kann. Im Falle eines Widerrufs wird der ensprechende Titel und Name aus der Datenbank des Campus Navigators gelöscht und ist somit nicht mehr weltweit abrufbar.</p>
                        </form>
                        
                        <h2>Raumbelegungsplan</h2>
                        
                        <p style='color:#ff0000'>Die Raumbelegungspläne dienen der unverbindlichen Vorinformation. Für verbindliche Buchungen wenden Sie sich bitte an die <a href='https://tu-dresden.de/service/arbeiten_tud/raumvermietung/' target='_blank'>operative Raumvergabe</a>.</p>
                        <p><a href="raum/542100.2130/drucken" class='druckansicht right'>Druckansicht</a></p>
                        <p><span class='bld'>Aktuelle Woche (Woche vom 21.11.2022 - 27.11.2022 47.KW)</span></p>
                        <table style='width:100%; max-width:812px; id='woche1'><tr><th style='width:50px;'>Uhrzeit</th><th>Montag</th><th>Dienstag</th><th>Mittwoch</th><th>Donnerstag</th><th>Freitag</th></tr><tr><td class='cent'>7:30 - 9:00</td><td></td><td></td><td></td><td><div>U Rechnerarchitektur I</div><span class='sml'><b>N.N.40 RA</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Rechnerarchitektur I</span></td><td><div>U SoI / HSC</div><span class='sml'><b>Wollschlaeger</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Systemorientierte Informatik und Hardware Software Codesign</span></td></tr><tr><td class='cent'>9:20 - 10:50</td><td><div>U Complexity Theory</div><span class='sml'><b>Krötzsch</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Complexity Theory</span></td><td><div>U HPC</div><span class='sml'><b>Nagel</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>High Performance Computing</span></td><td><div>V Principl. Dep. Syst.</div><span class='sml'><b>Fetzer</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Principles of Dependable Systems</span></td><td><div>U Informat.I/ ET</div><span class='sml'><b>N.N.10 ADS</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Informatik I für ET/MT/RES</span></td><td></td></tr><tr><td class='cent'>11:10 - 12:40</td><td><div>V Inf-Anw. Automation</div><span class='sml'><b>Wollschlaeger</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Informatik-Anwendungen in der Automation</span></td><td><div>U Betriebssysteme u. Sich.</div><span class='sml'><b>N.N.14 BSS</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Betriebssysteme und Sicherheit</span></td><td></td><td><div>U DB-Eng. Ü</div><span class='sml'><b>Lehner</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Datenbank-Engineering Übung </span></td><td><div>V Eng. Adapt. Mobile Apps</div><span class='sml'><b>Springer</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Engineering Adaptive Mobile Applications</span></td></tr><tr><td class='cent'>13:00 - 14:30</td><td><div>V Netzw. ind. Anw.</div><span class='sml'><b>Wollschlaeger</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Netzwerkmanagement in industriellen Anwendungen</span></td><td><div>U HS Techn. Datensch.</div><span class='sml'><b>Köpsell</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Hauptseminar Technischer Datenschutz</span></td><td><div>V PAofCS</div><span class='sml'><b>Nagel</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Performance Analysis of Computing Systems</span></td><td><div>U Principl. Dep. Syst.</div><span class='sml'><b>Fetzer</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Principles of Dependable Systems</span></td><td><div>U Eng. Adapt. Mobile Apps</div><span class='sml'><b>Springer</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Engineering Adaptive Mobile Applications</span></td></tr><tr><td class='cent'>14:50 - 16:20</td><td><div>V Information Retrieval</div><span class='sml'><b>Lehner</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Information Retrieval</span></td><td><div>U Microk.bas.</div><span class='sml'><b>Roitzsch</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Mikrokernbasierte Betriebssysteme</span></td><td><div>U Computergestützte Chirurgie</div><span class='sml'><b>Speidel</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Computer- und robotergestütze Chirurgie</span></td><td><div>U Form.Syst.</div><span class='sml'><b>N.N.33 FS</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Formale Systeme</span></td><td><div>U Informat.I/ ET</div><span class='sml'><b>N.N.13 ADS</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Informatik I für ET/MT/RES</span></td></tr><tr><td class='cent'>16:40 - 18:10</td><td><div>U Algorit.u.Daten</div><span class='sml'><b>N.N.03 AuD</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Algorithmen und Datenstrukturen</span></td><td><div>U Rechnerarchitektur I</div><span class='sml'><b>N.N.39 RA</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Rechnerarchitektur I</span></td><td><div>U SWT2</div><span class='sml'><b>N.N.17</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Softwaretechnologie II</span></td><td><div>U Info5 AuD</div><span class='sml'><b>N.N.INF</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Info5 AuD</span></td><td></td></tr><tr><td class='cent'>18:30 - 20:00</td><td></td><td></td><td></td><td><div>U Info5 AuD</div><span class='sml'><b>N.N.INF</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Info5 AuD</span></td><td></td></tr><tr><td class='cent'>20:20 - 21:50</td><td></td><td></td><td></td><td></td><td></td></tr><tr><td class='cent'>22:10 - 23:40</td><td></td><td></td><td></td><td></td><td></td></tr></table><p>Quelle: Datenbank Wintersemester</p>
                        <p style="margin-top: 16px;"><span class='bld'>Nächste Woche (Woche vom 28.11.2022 - 04.12.2022 48.KW)</span></p>
                        <table style='width:100%; max-width:812px; id='woche2'><tr><th style='width:50px;'>Uhrzeit</th><th>Montag</th><th>Dienstag</th><th>Mittwoch</th><th>Donnerstag</th><th>Freitag</th></tr><tr><td class='cent'>7:30 - 9:00</td><td></td><td><div>U Einf. Mediengest.</div><span class='sml'><b>Groh</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Einführung in die Mediengestaltung</span></td><td></td><td><div>U Rechnerarchitektur I</div><span class='sml'><b>N.N.40 RA</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Rechnerarchitektur I</span></td><td><div>U SoI / HSC</div><span class='sml'><b>Wollschlaeger</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Systemorientierte Informatik und Hardware Software Codesign</span></td></tr><tr><td class='cent'>9:20 - 10:50</td><td><div>U Complexity Theory</div><span class='sml'><b>Krötzsch</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Complexity Theory</span></td><td><div>U HPC</div><span class='sml'><b>Nagel</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>High Performance Computing</span></td><td><div>V Principl. Dep. Syst.</div><span class='sml'><b>Fetzer</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Principles of Dependable Systems</span></td><td><div>U Informat.I/ ET</div><span class='sml'><b>N.N.12 ADS</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Informatik I für ET/MT/RES</span></td><td></td></tr><tr><td class='cent'>11:10 - 12:40</td><td><div>V Inf-Anw. Automation</div><span class='sml'><b>Wollschlaeger</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Informatik-Anwendungen in der Automation</span></td><td><div>U Betriebssysteme u. Sich.</div><span class='sml'><b>N.N.14 BSS</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Betriebssysteme und Sicherheit</span></td><td></td><td><div>U DB-Eng. Ü</div><span class='sml'><b>Lehner</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Datenbank-Engineering Übung </span></td><td><div>V Eng. Adapt. Mobile Apps</div><span class='sml'><b>Springer</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Engineering Adaptive Mobile Applications</span></td></tr><tr><td class='cent'>13:00 - 14:30</td><td><div>V Netzw. ind. Anw.</div><span class='sml'><b>Wollschlaeger</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Netzwerkmanagement in industriellen Anwendungen</span></td><td><div>U HS Techn. Datensch.</div><span class='sml'><b>Köpsell</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Hauptseminar Technischer Datenschutz</span></td><td><div>V PAofCS</div><span class='sml'><b>Nagel</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Performance Analysis of Computing Systems</span></td><td><div>U Principl. Dep. Syst.</div><span class='sml'><b>Fetzer</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Principles of Dependable Systems</span></td><td><div>U Eng. Adapt. Mobile Apps</div><span class='sml'><b>Springer</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Engineering Adaptive Mobile Applications</span></td></tr><tr><td class='cent'>14:50 - 16:20</td><td><div>V Information Retrieval</div><span class='sml'><b>Lehner</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Information Retrieval</span></td><td><div>P KP Mikrokern basierter BS</div><span class='sml'><b>Roitzsch</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Komplexpraktikum Mikrokernbasierte Betriebssysteme</span></td><td><div>U Computergestützte Chirurgie</div><span class='sml'><b>Speidel</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Computer- und robotergestütze Chirurgie</span></td><td><div>U Form.Syst.</div><span class='sml'><b>N.N.33 FS</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Formale Systeme</span></td><td></td></tr><tr><td class='cent'>16:40 - 18:10</td><td><div>U Algorit.u.Daten</div><span class='sml'><b>N.N.03 AuD</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Algorithmen und Datenstrukturen</span></td><td><div>U Rechnerarchitektur I</div><span class='sml'><b>N.N.39 RA</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Rechnerarchitektur I</span></td><td><div>U SWT2</div><span class='sml'><b>N.N.17</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Softwaretechnologie II</span></td><td><div>U Info5 AuD</div><span class='sml'><b>N.N.INF</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Info5 AuD</span></td><td></td></tr><tr><td class='cent'>18:30 - 20:00</td><td></td><td></td><td></td><td><div>U Info5 AuD</div><span class='sml'><b>N.N.INF</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Info5 AuD</span></td><td></td></tr><tr><td class='cent'>20:20 - 21:50</td><td></td><td></td><td></td><td></td><td></td></tr><tr><td class='cent'>22:10 - 23:40</td><td></td><td></td><td></td><td></td><td></td></tr></table><p>Quelle: Datenbank Wintersemester</p>
                        
                        <div id="barf_div" style="display:none;">
                            <h2>Informationen zur Barrierefreiheit des Raumes</h2>
                            <table id="barf_info"></table>
                            <p id="barf_info_effective"></p>
                            <div id="editbox" style="display:none;">
                                <div class="closecontainer">
                                    <button id="closebutton" class="close" aria-label="Close editbox" type="button" style="background-color: #ddd;right:10px;top:-45px;">
                                        <span aria-hidden="true">&times;</span>
                                    </button>
                                </div>
                                <label for="value_corrected">Korrigierter Wert:</label><br>
                                <input id="value_corrected" type="text"><br>
                                <label for="feedbackmail">Mailadresse für Feedback (optional):</label><br>
                                <input id="feedbackmail" type="email">
                                <input id="send_corrected_value" type="button" value="OK" style="float:right;cursor:pointer;">
                            </div>
                        </div>
                        
                        <div id="routing_buttons" style="display: none;">
                            <h2>Routing Start/Ziel</h2>
                            <button class="big" onclick="location.href = 'routing/apb E008/-/foot,shortest,blind/';"
                                    style="height:40px;font-weight:bold;background: #002557;padding:0px 5px 0px 5px;color:#ffffff;border:1px solid #cccccc;cursor:pointer;">
                                Routing von diesem Raum
                            </button>
                            <button class="big" onclick="location.href = 'routing/-/apb E008/foot,shortest,blind/';"
                                    style="height:40px;font-weight:bold;background: #002557;padding:0px 5px 0px 5px;color:#ffffff;border:1px solid #cccccc;cursor:pointer;">
                                Routing zu diesem Raum
                            </button>
                        </div>
                        <h2>Links</h2>
                        <ul id="link_list" class="reg">
                            
                            <li class="reg">
                                <span id="input_view_in_floor">Anzeige im Etagenplan:</span> 
                                <input aria-labelledby="input_view_in_floor" type='text' value='https://navigator.tu-dresden.de/etplan/apb/00/raum/542100.2130' readonly='readonly' onclick='this.select()' />
                            </li>
                            <li class="reg">
                                <span id="input_this_site">Diese Seite:</span> 
                                <input aria-labelledby="input_this_site" type='text' value='https://navigator.tu-dresden.de/raum/542100.2130' readonly='readonly' onclick='this.select()' />
                            </li>
                        </ul>
                    </div>
                    <div id="mouseinfo" style="position: absolute; background:#687372; display: none; padding: 4px; color:#ffffff; font-size:10px; line-height: 14px; opacity: 0.85; border: 1px solid #ffffff; left:0px; top:0px;">
                    </div>
                </div>
                <div class="clear">
                </div>
            </div>
        </div>
    </body>
    <footer id="footer">
        




<div id="privacy-statement">
    <span>Unsere Webseite nutzt Cookies und die Analysesoftware Matomo.
        <input id="custom_cookies" class="checkboxwithlabel" type="checkbox" style="margin-right:0px;">
        <label for="custom_cookies">Personalisierungscookies</label>
        <input type="button" value="Ok" onclick="cn.Main.hidePrivStmt();"/>
    </span>
    <p class="sml">Es werden Cookies, die zur Erbringung der Dienstleistung zwingend erforderlich sind, genutzt. Das heißt insbesondere werden keine sogenannten Tracking-Cookies genutzt, um Nutzerbewegungen und das Surfverhalten der Nutzenden unserer Seite zu erfassen bzw. zu analysieren. Darüberhinaus können die Nutzenden der Speicherung von Cookies zustimmen die die Personalisierung der Webseite ermöglichen (Routingeinstellungen, Vollbild etc.).</p>
</div>




<div id="portal-footer" >
    <a title="campus-navigator[a_t]mailbox[do_t]tu-dresden[do_t]de" class="asm">Kontakt</a>
    <a href="impressum">Impressum</a>
    <a href="impressum#datenschutz">Datenschutz</a>
    <a href="barrierenwebseite">Barrierefreiheit</a>
    <a href="http://tu-dresden.de"><img id="tu-logo" src='images/logo.png' alt="Logo der TU Dresden" title="Startseite der Technischen Universität Dresden" /></a>
	<span style="color:white">Datenstand: 29.09.2022<span>
</div>


    </footer>
</html>
`
}
