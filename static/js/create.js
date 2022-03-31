
var c=0; //счётчик количества строк
function addline()
{
  c++; // увеличиваем счётчик строк
  s=document.getElementById('table').innerHTML; // получаем HTML-код таблицы
  s=s.replace(/[\r\n]/g,''); // вырезаем все символы перевода строк
  re=/(.*)(<tr id=.*>)(<\/table>)/gi;
                // это регулярное выражение позволяет выделить последнюю строку таблицы
  s1=s.replace(re,'$2'); // получаем HTML-код последней строки таблицы
  s2=s1.replace(/\[\d+\]/gi,'['+c+']'); // заменяем все цифры к квадратных скобках
                // на номер новой строки
  s2=s2.replace(/(rmline\()(\d+\))/gi,'$1'+c+')');
                // заменяем аргумент функции rmline на номер новой строки
  s=s.replace(re,'$1$2'+s2+'$3');
                // создаём HTML-код с добавленным кодом новой строки
  document.getElementById('table').innerHTML=s;
                // возвращаем результат на место исходной таблицы
  return false; // чтобы не происходил переход по ссылке
}
function rmline(q)
{
                 if (c==0) return false; else c--;
                // если раскомментировать предыдущую строчку, то последний (единственный)
                // элемент удалить будет нельзя.
  s=document.getElementById('table').innerHTML;
  s=s.replace(/[\r\n]/g,'');
  re=new RegExp('<tr id="?newline"? nomer="?\\['+q+'.*?<\\/tr>','gi');
                // это регулярное выражение позволяет выделить строку таблицы с заданным номером
  s=s.replace(re,'');
                // заменяем её на пустое место
  document.getElementById('table').innerHTML=s;
  return false;
}


var v=0; //счётчик количества строк для дополнительных предложений
function addlineP()
{
  v++; // увеличиваем счётчик строк
  s=document.getElementById('table_1').innerHTML; // получаем HTML-код таблицы
  s=s.replace(/[\r\n]/g,''); // вырезаем все символы перевода строк
  re=/(.*)(<tr id=.*>)(<\/table>)/gi;
                // это регулярное выражение позволяет выделить последнюю строку таблицы
  s1=s.replace(re,'$2'); // получаем HTML-код последней строки таблицы
  s2=s1.replace(/\[\d+\]/gi,'['+v+']'); // заменяем все цифры к квадратных скобках
                // на номер новой строки
  s2=s2.replace(/(rmlineP\()(\d+\))/gi,'$1'+v+')');
                // заменяем аргумент функции rmline на номер новой строки
  s=s.replace(re,'$1$2'+s2+'$3');
                // создаём HTML-код с добавленным кодом новой строки
  document.getElementById('table_1').innerHTML=s;
                // возвращаем результат на место исходной таблицы
  return false; // чтобы не происходил переход по ссылке
}
function rmlineP(z)
{
                 if (v==0) return false; else v--;
                // если раскомментировать предыдущую строчку, то последний (единственный)
                // элемент удалить будет нельзя.
  s=document.getElementById('table_1').innerHTML;
  s=s.replace(/[\r\n]/g,'');
  re=new RegExp('<tr id="?newlineP"? nomer="?\\['+z+'.*?<\\/tr>','gi');
                // это регулярное выражение позволяет выделить строку таблицы с заданным номером
  s=s.replace(re,'');
                // заменяем её на пустое место
  document.getElementById('table_1').innerHTML=s;
  return false;
}



// Example starter JavaScript for disabling form submissions if there are invalid fields
// (function () {
//   'use strict'
//
//   // Fetch all the forms we want to apply custom Bootstrap validation styles to
//   var forms = document.querySelectorAll('.needs-validation')
//
//   // Loop over them and prevent submission
//   Array.prototype.slice.call(forms)
//     .forEach(function (form) {
//       form.addEventListener('submit', function (event) {
//         if (!form.checkValidity()) {
//           event.preventDefault()
//           event.stopPropagation()
//         }
//
//         form.classList.add('was-validated')
//       }, false)
//     })
// })()

//проверка размера файла

var uploadField = document.getElementById("file");

uploadField.onchange = function() {
    if(this.files[0].size > 2048576){
       alert("File is too big!");
       this.value = "";
    };
};

