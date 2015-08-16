var Contact = function () {
  var dst = window.location.protocol + "//" + window.location.hostname + '/contact.php';
  function responseSuccess(data) {
    data = JSON.parse(data);
    if(data.status === 'success') {
      $('#contactform').html('Submission sent succesfully.');
    } else {
      $('#contactform').html('Submission failed, please contact directly.');
    }
  }
  return {
    init: function () {
      $('#contactform').validator()
      $('#contactform').validator().on('submit',function(e) {
       if (e.isDefaultPrevented()){

       }else{
          $('#contactform *').fadeOut(200);
          $('#contactform').prepend('Your submission is being processed...');

          $.ajax({
            type     : 'POST',
            cache    : false,
            url      : dst,
            data     : $(this).serialize(),
            success  : function(data) {
              responseSuccess(data);
            },
            error  : function(data) {
            }
          });
          return false;
        }

      });
    }
  }
}();
