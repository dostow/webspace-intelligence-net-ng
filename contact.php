<?php

if(empty($_POST) || !isset($_POST)) {

  ajaxResponse('error', 'Post cannot be empty.');

} else {

  $postData = $_POST;
  $dataString = implode($postData,",");

  $mailgun = sendMailgun($postData);

  if($mailgun) {
    ajaxResponse('success', 'Great success.', $postData, $mailgun);
  } else {
    ajaxResponse('error', 'Mailgun did not connect properly.', $postData, $mailgun);
  }

}

function ajaxResponse($status, $message, $data = NULL, $mg = NULL) {
  $response = array (
    'status' => $status,
    'message' => $message,
    'data' => $data,
    'mailgun' => $mg
    );
  $output = json_encode($response);
  exit($output);
}

function sendMailgun($data) {

  $api_key = 'key-aae2b80f42abb6fe502c5ff9fea8d31e';
  $api_domain = 'idcgroupltd.com';
  $send_to = 'idcgroupltd@gmail.com';

  // sumbission data
  $ipaddress = $_SERVER['REMOTE_ADDR'];
  $date = date('d/m/Y');
  $time = date('H:i:s');

  // form data
  $name = $data['name'];
  $postcontent = $data['message'];
  $reply = $data['email'];

  $messageBody = "<p>You have received a new message from {$name}</p>
                {$postcontent}
                <p>This message was sent from the IP Address: {$ipaddress} on {$date} at {$time}</p>";

  $config = array();
  $config['api_key'] = $api_key;
  $config['api_url'] = 'https://api.mailgun.net/v3/'.$api_domain.'/messages';

  $message = array();
  $message['from'] = $reply;
  $message['to'] = $send_to;
  $message['h:Reply-To'] = $reply;
  $message['subject'] = "New Message from intelligence.net.ng Contact Form";
  $message['html'] = $messageBody;

  $curl = curl_init();

  curl_setopt($curl, CURLOPT_URL, $config['api_url']);
  curl_setopt($curl, CURLOPT_HTTPAUTH, CURLAUTH_BASIC);
  curl_setopt($curl, CURLOPT_USERPWD, "api:{$config['api_key']}");
  curl_setopt($curl, CURLOPT_RETURNTRANSFER, 1);
  curl_setopt($curl, CURLOPT_CONNECTTIMEOUT, 10);
  curl_setopt($curl, CURLOPT_SSL_VERIFYPEER, 0);
  curl_setopt($curl, CURLOPT_SSL_VERIFYHOST, 0);
  curl_setopt($curl, CURLOPT_POST, true);
  curl_setopt($curl, CURLOPT_POSTFIELDS,$message);

  $result = curl_exec($curl);
  curl_close($curl);
  return $result;

}
 ?>
