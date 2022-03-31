// Закрытие аккордиона в меню при закрытии меню (на телефоне)

//(function () {
//    var flush_collapseReadyHolidays = document.getElementById('flush-collapseReadyHolidays');
//    //var accordion_button = document.querySelector('.accordion-button');
//    var burger_menu = document.getElementById('burger_menu');
//    if(burger_menu.classList.contains('burger_nav_active')) {

//    } else {
//      //accordion_button.classList.add('collapse');
//      flush_collapseReadyHolidays.classList.remove('show');
//accordion_button.setAttribute("aria-expanded", "false");
//    }
//  }());





//Отправка формы на почту form_down
//(
//"use strict"

//document.addEventListener('DOMContentLoaded', function () {
//  const form_down = document.getElementById('form_down');
//  const spinner = document.querySelector('.my_spinner_down');
//  form_down.addEventListener('submit', formSend);

//  async function formSend(e) {
//    e.preventDefault();
//    let error = formValidate(form_down);

//  var form_downData = new FormData(form_down);
//  console.log(error);

//  if (error === 0) {
//    spinner.classList.add('_sending');
//
//    let response = await fetch('sendmail.php', {
//      method: 'POST'
//    });
//    if (response.ok) {
//    let result = await response.json();
//      alert(result.message);
//      form_down.reset();
//      spinner.classList.remove('_sending');
//    } else {
//      alert('Ошибка');
//      spinner.classList.remove('_sending');
//    }
//} else {
//    alert('Заполните обязательные поля') //можно сделать popup
//  }
//  }

//function formValidate(form_down){
//  let error = 0;
//  let formReq = document.querySelectorAll('._req');

//  for (let index = 0; index < formReq.length; index++) {
//  const input = formReq[index];
//  formRemoveError(input);

//  if (input.classList.contains('_tel')){
//    if ((input.value === '+7') || (input.value === '')) {
//        formAddError(input);
//        error++;
//    }
//    }
//    if (input.classList.contains('_name')){
//      if (input.value === '') {
//        formAddError(input);
//        error++;
//      }
//      }
//    }
//    return error;
//}
//)

//Проверка формы на заполнение полей

function sendform() {

    if (document.forms[0].tel.value == "") { //|| document.forms[0].tel.value == "+7") {
        alert('Пожалуйста, введите телефон');
        document.mailform.tel.focus();
        return false
    }
    if (document.forms[0].name.value == "") {
        alert('Пожалуйста, введите Ваше имя');
        document.mailform.name.focus();
        return false
    }
    return true;
}

function validateForm() {
    var x = document.forms["name"]["tel"].value;
    if (x == "") {
        alert("Name must be filled out");
        return false;
    }
    return true;
}

// window.addEventListener('resize',(e)=> { 
//   document.querySelector('.submenu').style.transform = 
//     document.body.clientWidth>1140 ? '' : 
//       `scale(${document.body.clientWidth/1380}) translate(${-(1140-document.body.clientWidth)/2}px)`; 
// });

// window.addEventListener('resize', (e) => {
//     document.querySelector('.submenu').style.transform =
//         document.body.clientWidth > 1210 ? '' :
//         `translate(${-(1210-document.body.clientWidth)/2}px)`;
// });

//центрирование подменю относительно пунктов меню
// let topmenu;
// let submenu;
// document.onmouseover = function(e) {


//     let topmenu = e.target.closest('.topmenu_item')
//     if (!topmenu) return
//     let submenu = topmenu.querySelector('.submenu')
//     if (!submenu) return
//     let coords = topmenu.getBoundingClientRect()

//     let left = coords.left + (topmenu.offsetWidth - submenu.offsetWidth) / 2 + submenu.offsetWidth / 2

//     if (left < submenu.offsetWidth / 2) left = submenu.offsetWidth / 2 + 5
//     submenu.style.left = left + 'px'
//     submenu.style.visibility = 'visible'

// };
// document.onmouseout = function(e) {
//     if (submenu) {
//         topmenu = null
//         submenu = null
//     }
// };

let topmenu;
let submenu;
document.onmouseover = function(e) {


    let topmenu = e.target.closest('.topmenu_item')
    if (!topmenu) return
    let submenu = topmenu.querySelector('.submenu')
    if (!submenu) return
    let coords = topmenu.getBoundingClientRect()

    let left = coords.left + (topmenu.offsetWidth - submenu.offsetWidth) / 2 + submenu.offsetWidth / 2

    if (left < submenu.offsetWidth / 2) left = submenu.offsetWidth / 2 + 5
    submenu.style.left = left + 'px'
};
document.onmouseout = function(e) {
    if (submenu) {
        topmenu = null
        submenu = null
    }
};


// window.addEventListener('resize',(e)=> { 
//   document.querySelector('.submenu').style.transform = 
//     document.body.clientWidth>1200 ? '' : 
//       ` translate(${-(1200-document.body.clientWidth)/2}px)`; 
// });


// document.body.querySelector('.submenu').onmouseout = function(e) {
//     let submenu = e.target.closest('.submenu')
//     submenu.style.left = 0 + 'px'
// };

//Функия для проверки полей
//function formAddError(input) {
//    input.parentElement.classList.add('_error');
//  input.classList.add('_error');
//  }
//function formRemoveError(input) {
//  input.parentElement.classList.remove('_error');
//  input.classList.remove('_error');
//  }
//Функция теста номера телефона
//  function telTest(input) {
//  return /^[+78]+$/.test(input.value);
//  }
//});

//Замена карусели в зависимости от размера экрана

(function() {
    const carouselStockFullScreen = document.getElementById('carouselStockFullScreen');
    const carouselStockMobileScreen = document.getElementById('carouselStockMobileScreen');
    if (window.innerWidth <= 768 && carouselStockFullScreen != null) {
        carouselStockFullScreen.classList.add('hidden_all');
        carouselStockMobileScreen.classList.remove('hidden_all');
    } else if (carouselStockFullScreen != null) {
        carouselStockFullScreen.classList.remove('hidden_all');
        carouselStockMobileScreen.classList.add('hidden_all');
    }
}());


// Форма обратной связи появляющаяся внизу экрана при скроле и фиксированное меню

(function() {
    // для горищзонтального меню
    const header_menu = document.querySelector('.header_menu');
    const header = document.querySelector('.header');
    //для формы обратной связи

    const form = document.querySelector('.form_down');
    const upbutton = document.querySelector('.button_back');
    var pageHeight = document.documentElement.scrollHeight;
    var screenHeight = document.documentElement.clientHeight;
    var pageHeight_1 = pageHeight - (screenHeight + 350);
    //var pageHeight_1 = pageHeight - (screenHeight + 250);

    window.onscroll = () => {
        // для горизонтального меню
        if (window.pageYOffset >= 90 && window.innerWidth > 950) {
            header_menu.classList.add('menu_active');
            header.classList.add('header_active');
        } else {
            header_menu.classList.remove('menu_active');
            header.classList.remove('header_active');
        }

        //для формы обратной связи внизу экрана
        if (window.pageYOffset > pageHeight_1) {
            form.classList.add('hidden_form');
            form.classList.remove('visible_form');
        } else if (window.pageYOffset > 150 && window.innerWidth >= 1024) { //Редактировать высоту экрана в пикселях!!!
            form.classList.remove('hidden_form');
            form.classList.add('visible_form');
        } else {
            form.classList.add('hidden_form');
            form.classList.remove('visible_form');
        }
        //появление upbotton
        if (window.pageYOffset > 500 && window.innerWidth >= 768) {
            upbutton.classList.add('upbutton_visiable');
            upbutton.classList.remove('upbutton_hidden');
        } else {
            upbutton.classList.add('upbutton_hidden');
            upbutton.classList.remove('upbutton_visiable');
        }
    };
}());


//Маска для номера телефона


//$(document).ready(function() {
// $('._tel').mask('+7(000)000-00-00');
//});

//$(function(){
//2. Получить элемент, к которому необходимо добавить маску
//$("#tel_down").mask("8(999) 999-9999");
//});

window.addEventListener("DOMContentLoaded", function() {
    [].forEach.call(document.querySelectorAll('._tel'), function(input) {
        var keyCode;

        function mask(event) {
            event.keyCode && (keyCode = event.keyCode);
            var pos = this.selectionStart;
            if (pos < 3) event.preventDefault();
            var matrix = "+7 (___) ___ ____",
                i = 0,
                def = matrix.replace(/\D/g, ""),
                val = this.value.replace(/\D/g, ""),
                new_value = matrix.replace(/[_\d]/g, function(a) {
                    return i < val.length ? val.charAt(i++) || def.charAt(i) : a
                });
            i = new_value.indexOf("_");
            if (i != -1) {
                i < 5 && (i = 3);
                new_value = new_value.slice(0, i)
            }
            var reg = matrix.substr(0, this.value.length).replace(/_+/g,
                function(a) {
                    return "\\d{1," + a.length + "}"
                }).replace(/[+()]/g, "\\$&");
            reg = new RegExp("^" + reg + "$");
            if (!reg.test(this.value) || this.value.length < 5 || keyCode > 47 && keyCode < 58) this.value = new_value;
            if (event.type == "blur" && this.value.length < 5) this.value = ""
        }

        input.addEventListener("input", mask, false);
        input.addEventListener("focus", mask, false);
        input.addEventListener("blur", mask, false);
        input.addEventListener("keydown", mask, false)

    });

});



//Прелдоадер
window.onload = function() {
    document.body.classList.add('loaded_hiding');
    window.setTimeout(function() {
        document.body.classList.add('loaded');
        document.body.classList.remove('loaded_hiding');
    }, 500);

    if ((+localStorage.getItem('utm_timeout') + 86400000) < Date.now()) {
        localStorage.setItem('source', "");
        localStorage.setItem('medium', "");
        localStorage.setItem('campaign', "");
        localStorage.setItem('content', "");
        localStorage.setItem('term', "");
    }

    if (location.search != "") {
        // Parse the URL
        function getParameterByName(name) {
            var name = name.replace(/[\[]/, "\\[").replace(/[\]]/, "\\]");
            var regex = new RegExp("[\\?&]" + name + "=([^&#]*)");
            var results = regex.exec(location.search);
            return results === null ? "" : decodeURIComponent(results[1].replace(/\+/g, " "));
        }
        // Give the URL parameters variable names
        var source = getParameterByName('utm_source');
        var medium = getParameterByName('utm_medium');
        var campaign = getParameterByName('utm_campaign');
        var content = getParameterByName('utm_content');
        var term = getParameterByName('utm_term');

        localStorage.setItem('source', source);
        localStorage.setItem('medium', medium);
        localStorage.setItem('campaign', campaign);
        localStorage.setItem('content', content);
        localStorage.setItem('term', term);
        localStorage.setItem('utm_timeout', Date.now());
        //U Can use it with Expertsender Forms
    }

    //Put the variable names into the hidden fields in the form.

    document.getElementById("property_utm_source_2").value = localStorage.getItem('source');
    document.getElementById("property_utm_medium_2").value = localStorage.getItem('medium');
    document.getElementById("property_utm_campaign_2").value = localStorage.getItem('campaign');
    document.getElementById("property_utm_content_2").value = localStorage.getItem('content');
    document.getElementById("property_utm_term_2").value = localStorage.getItem('term');
    document.getElementById("property_utm_source_3").value = localStorage.getItem('source');
    document.getElementById("property_utm_medium_3").value = localStorage.getItem('medium');
    document.getElementById("property_utm_campaign_3").value = localStorage.getItem('campaign');
    document.getElementById("property_utm_content_3").value = localStorage.getItem('content');
    document.getElementById("property_utm_term_3").value = localStorage.getItem('term');
    document.getElementById("property_utm_source_1").value = localStorage.getItem('source');
    document.getElementById("property_utm_medium_1").value = localStorage.getItem('medium');
    document.getElementById("property_utm_campaign_1").value = localStorage.getItem('campaign');
    document.getElementById("property_utm_content_1").value = localStorage.getItem('content');
    document.getElementById("property_utm_term_1").value = localStorage.getItem('term');
    
}

// Появление кнопки отправки при установке галочки в чекбоксе

function fun_check() {
    var chbox = document.getElementById('check_form_intro');
    var btn = document.getElementById('button_form_intro');
    const button_1 = document.querySelector('.search-form_submit');

    if (!chbox.checked) {
        btn.setAttribute('disabled', true);
        button_1.classList.remove('search-form_submit_active');

    } else {
        btn.removeAttribute('disabled');
        button_1.classList.add('search-form_submit_active');

    }
}

// Появление кнопки отправки при установке галочки в чекбоксе для модального окна

function fun_check_modal() {
    var chbox = document.getElementById('check_form_modal');
    var btn = document.getElementById('modal_button_form');
    //const button_1 = document.querySelector('.modal_search-form_submit');

    if (!chbox.checked) {
        btn.setAttribute('disabled', true);
        btn.classList.remove('search-form_submit_active_modal');
    } else {
        btn.removeAttribute('disabled');
        btn.classList.add('search-form_submit_active_modal');
    }
}

// Появление кнопки отправки при установке галочкb в чекбоксе

function fun_check_down() {
    var chbox = document.getElementById('check_form_down');
    var btn = document.getElementById('button_form_down');
    

    if (!chbox.checked) {
        btn.setAttribute('disabled', true);
        btn.classList.remove('search-form_submit_active_benefit');

    } else {
        btn.removeAttribute('disabled');
        btn.classList.add('search-form_submit_active_benefit');
    }
}

// Появление кнопки отправки при установке галочкb в чекбоксе

function fun_check_benefit() {
    var chbox = document.getElementById('check_form_benefit');
    var btn = document.getElementById('benefit_button_form');
    

    if (!chbox.checked) {
        btn.setAttribute('disabled', true);
        btn.classList.remove('search-form_submit_active_benefit');
    } else {
        btn.removeAttribute('disabled');
        btn.classList.add('search-form_submit_active_benefit');
    }
}

//Burger handler

(function() {
    const burgerItem = document.querySelector('.burger');
    const menu = document.getElementById('burger_menu');
    const body = document.body;
    const menuLinks = document.querySelectorAll('.menu_link_burger');
    const menuLinksSub = document.querySelectorAll('.menu_link_sub_burger');
    burgerItem.addEventListener('click', () => {
        menu.classList.toggle('burger_nav_active');
        burgerItem.classList.toggle('burger_active');
        body.classList.toggle('stop-scrolling');
    });
    if (window.innerWidth <= 950) {
        for (let i = 0; i < menuLinks.length; i++) {
            menuLinks[i].addEventListener('click', () => {
                menu.classList.remove('burger_nav_active');
                body.classList.toggle('stop-scrolling');
                burgerItem.classList.toggle('burger_active');
            });
        }
        for (let i = 0; i < menuLinksSub.length; i++) {
            menuLinksSub[i].addEventListener('click', () => {
                menu.classList.remove('burger_nav_active');
                body.classList.toggle('stop-scrolling');
                burgerItem.classList.toggle('burger_active');
            });
        }
    }
}());

//Scroll

// Scroll to anchors
document.querySelectorAll('a[href^="#"').forEach(link => {

    link.addEventListener('click', function(e) {
        e.preventDefault();

        let href = this.getAttribute('href').substring(1);

        const scrollTarget = document.getElementById(href);

        const topOffset = document.querySelector('.header').offsetHeight;
        // const topOffset = 0; // если не нужен отступ сверху
        const elementPosition = scrollTarget.getBoundingClientRect().top;
        const offsetPosition = elementPosition - topOffset;

        window.scrollBy({
            top: offsetPosition,
            behavior: 'smooth'
        });
    });
});


// Назначение стилей активному пункту меню
// (function() {
//     let currentLi;
//     let menu = document.querySelector('.topmenu')
//     menu.onclick = function(e) {

//         let li = e.target.closest('li')
//         if (!li) return
//         if (!menu.contains(li)) return

//         hover(li)
//     }

//     function hover(li) {
//         if (currentLi) {
//             currentLi.classList.remove('active_header_menu')
//         }
//         currentLi = li
//         currentLi.classList.add('active_header_menu')
//     }
// }());

// var menuItems = document.querySelector('.topmenu')
// var onClick = function(event) {
//     event.preventDefault();

//     console.log(665)

//     for (var i = 0; i < menuItems.length; i++) {
//         menuItems[i].classList.remove('active_header_menu');
//     }

//     event.currentTarget.classList.add('active_header_menu');
// };

// for (var i = 0; i < menuItems.length; i++) {
//     menuItems[i].addEventListener('click', onClick, false);
// }

//Запрет выделения текста на сайте

// document.body.onmousedown = function(e) {
//     e.preventDefault()
// }

// выделение активного пункта меню
document.addEventListener('DOMContentLoaded', function() {
    let itemsMenu = document.querySelectorAll('.main_menu')
    let mainPoint = document.querySelector('.main')

    if (!mainPoint) return
    for (let item of itemsMenu) {
        if (item.dataset.chapter == mainPoint.dataset.chapter) {
            item.classList.add('menu_border')
        }
    }
})

//Уведомление об успешной отправки формы
// document.body.addEventListener('click', function(e) {
//     if (e.target.tagName != 'BUTTON') return
//     if (e.target.dataset.postform == 'true') {
//         localStorage.setItem('form_post', true)
//     }
// })

// document.addEventListener('DOMContentLoaded', function() {
//     if (localStorage.getItem('form_post') == true) {
//         document.getElementById('Modal_info').classList.remove('hidden_all')
//         localStorage.setItem('form_post', false)
//     }
// })


//Уведомление об успешной отправки формы

document.addEventListener('DOMContentLoaded', function() {
    let modalInfo = document.createElement('div')
    modalInfo.innerHTML = `

    `
    document.querySelector('.main').append(modalInfo)
})

// Уведомление об использовании файлов cookie

function checkCookies() {
    let cookieDate = localStorage.getItem('cookieDate');
    let cookieNotification = document.getElementById('cookie_notification');
    let cookieBtn = cookieNotification.querySelector('.cookie_accept');

    // Если записи про кукисы нет или она просрочена на 1 год, то показываем информацию про кукисы
    if (!cookieDate || (+cookieDate + 31536000000) < Date.now()) {
        cookieNotification.classList.add('show');
    }

    // При клике на кнопку, в локальное хранилище записывается текущая дата в системе UNIX
    cookieBtn.addEventListener('click', function() {
        localStorage.setItem('cookieDate', Date.now());
        cookieNotification.classList.remove('show');
    })
}
checkCookies();

// Передача в скрытые поля форм UTM меток

// Lazy download picture 
// document.addEventListener("DOMContentLoaded", function() {
//     var lazyloadImages = document.querySelectorAll("img.lazy");    
//     var lazyloadThrottleTimeout;
    
//     function lazyload () {
//       if(lazyloadThrottleTimeout) {
//         clearTimeout(lazyloadThrottleTimeout);
//       }   
      
//       lazyloadThrottleTimeout = setTimeout(function() {
//           var scrollTop = window.pageYOffset;
//           lazyloadImages.forEach(function(img) {
//               if(img.offsetTop < (window.innerHeight + scrollTop)) {
//                 img.src = img.dataset.src;
//                 img.classList.remove('lazy');
//               }
//           });
//           if(lazyloadImages.length == 0) { 
//             document.removeEventListener("scroll", lazyload);
//             window.removeEventListener("resize", lazyload);
//             window.removeEventListener("orientationChange", lazyload);
//           }
//       }, 20);
//     }
    
//     document.addEventListener("scroll", lazyload);
//     window.addEventListener("resize", lazyload);
//     window.addEventListener("orientationChange", lazyload);
//   });



