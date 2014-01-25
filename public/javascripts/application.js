$.getJSON('/index.json', function(data) {
  $('#busy').hide();
  
  if (data.error) { 
    alert(data.message);
  }
  else {
    var tbody = $('#gigs tbody');
    for (var i = 0; i < data.length; ++i) {
      tbody.append("<tr><td><a href='"+data[i]+"'>"+data[i]+"</a></td></tr>");
    }
  }
});
