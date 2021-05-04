var fs = require('fs')

const OPTION_A = process.env.OPTION_A;
const OPTION_B = process.env.OPTION_B;

fs.readFile("./views/index.html", 'utf8', function (err,data) {
  if (err) {
    return console.log(err);
  }
  var result = data.replace(/<div class="label">Cats<\/div>/g, `<div class="label">${OPTION_A}</div>`);
  var result = result.replace(/<div class="label">Dogs<\/div>/g, `<div class="label">${OPTION_B}</div>`);

  fs.writeFile("./views/index.html", result, 'utf8', function (err) {
     if (err) return console.log(err);
  });
});