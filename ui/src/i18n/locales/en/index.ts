import nlsCommon from "./nls.common.json";
import nlsLogin from "./nls.login.json";
import nlsDashboard from "./nls.dashboard.json";
import nlsSettings from "./nls.settings.json";
import nlsDomain from "./nls.domain.json";
import nlsAccess from "./nls.access.json";
import nlsWorkflow from "./nls.workflow.json";
import nlsCertificate from "./nls.certificate.json";

export default Object.freeze({
  ...nlsCommon,
  ...nlsLogin,
  ...nlsDashboard,
  ...nlsSettings,
  ...nlsDomain,
  ...nlsAccess,
  ...nlsWorkflow,
  ...nlsCertificate,
});
