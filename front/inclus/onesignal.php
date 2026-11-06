<?php
$appId = getenv('ONESIGNAL_APP_ID') ?: '';
?>
<script>
window.ONESIGNAL_APP_ID = <?= json_encode($appId, JSON_HEX_TAG | JSON_HEX_AMP | JSON_HEX_APOS | JSON_HEX_QUOT) ?>;
</script>
<script src="https://cdn.onesignal.com/sdks/web/v16/OneSignalSDK.page.js" defer></script>
<script src="/assets/js/onesignal.js" defer></script>
