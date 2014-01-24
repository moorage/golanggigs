$.getJSON('/index.json', function(data) {
  if (data.error) { 
    alert(data.message);
  }
  else {
    console.log(data);
  }
});
