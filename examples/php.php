<?
define('BASE_PATH', "/home/martero/workspace/php/tuenti-ng/");
ini_set('display_errors', 'On');
 
require_once BASE_PATH.'main/core/bootstrap/scriptsBootstrap.php';
$GLOBALS['request_time'] = time();
ClientInterface::setCurrentClientInterface(ClientInterface::INTERFACE_TYPE_API);
var_dump(TMainConfig::getConfigPaths());


use Tuenti\UserService\v2\constants\Realm;
$subscriptionsApi = SubscriptionsApi::create();
$msisdn = "541144374685";
$env = "preproduction";

$response = $subscriptionsApi->getFirstActiveByMsisdn($msisdn, Realm::MOVISTAR_AR);
//$response = $subscriptionsApi->getFirstActiveByMsisdn($msisdn, 2, $env);

var_dump($response);

