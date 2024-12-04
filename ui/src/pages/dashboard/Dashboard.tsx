import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { CalendarClock, CalendarX2, FolderCheck, SquareSigma, Workflow } from "lucide-react";

import { Statistic } from "@/domain/domain";

import { get } from "@/api/statistics";

import CertificateList from "@/components/certificate/CertificateList";

const Dashboard = () => {
  const [statistic, setStatistic] = useState<Statistic>();

  const { t } = useTranslation();

  useEffect(() => {
    const fetchStatistic = async () => {
      const data = await get();
      setStatistic(data);
    };

    fetchStatistic();
  }, []);

  return (
    <div className="flex flex-col">
      <div className="flex justify-between items-center">
        <div className="text-muted-foreground">{t("dashboard.page.title")}</div>
      </div>
      <div className="flex mt-10 gap-5 flex-col flex-wrap md:flex-row">
        <div className="w-full md:w-[250px] 3xl:w-[300px] flex items-center rounded-md p-3 shadow-lg border">
          <div className="p-3">
            <SquareSigma size={48} strokeWidth={1} className="text-blue-400" />
          </div>
          <div>
            <div className="text-muted-foreground font-semibold">{t("dashboard.statistics.all.certificate")}</div>
            <div className="flex items-baseline">
              <div className="text-3xl text-stone-700 dark:text-stone-200">
                {statistic?.certificateTotal ? (
                  <Link to="/certificate" className="hover:underline">
                    {statistic?.certificateTotal}
                  </Link>
                ) : (
                  0
                )}
              </div>
              <div className="ml-1 text-stone-700 dark:text-stone-200">{t("dashboard.statistics.unit")}</div>
            </div>
          </div>
        </div>

        <div className="w-full md:w-[250px] 3xl:w-[300px] flex items-center rounded-md p-3 shadow-lg border">
          <div className="p-3">
            <CalendarClock size={48} strokeWidth={1} className="text-yellow-400" />
          </div>
          <div>
            <div className="text-muted-foreground font-semibold">{t("dashboard.statistics.near_expired.certificate")}</div>
            <div className="flex items-baseline">
              <div className="text-3xl text-stone-700 dark:text-stone-200">
                {statistic?.certificateExpireSoon ? (
                  <Link to="/certificate?state=expireSoon" className="hover:underline">
                    {statistic?.certificateExpireSoon}
                  </Link>
                ) : (
                  0
                )}
              </div>
              <div className="ml-1 text-stone-700 dark:text-stone-200">{t("dashboard.statistics.unit")}</div>
            </div>
          </div>
        </div>

        <div className="border w-full md:w-[250px] 3xl:w-[300px] flex items-center rounded-md p-3 shadow-lg">
          <div className="p-3">
            <CalendarX2 size={48} strokeWidth={1} className="text-red-400" />
          </div>
          <div>
            <div className="text-muted-foreground font-semibold">{t("dashboard.statistics.expired.certificate")}</div>
            <div className="flex items-baseline">
              <div className="text-3xl text-stone-700 dark:text-stone-200">
                {statistic?.certificateExpired ? (
                  <Link to="/certificate?state=expired" className="hover:underline">
                    {statistic?.certificateExpired}
                  </Link>
                ) : (
                  0
                )}
              </div>
              <div className="ml-1 text-stone-700 dark:text-stone-200">{t("dashboard.statistics.unit")}</div>
            </div>
          </div>
        </div>

        <div className="border w-full md:w-[250px] 3xl:w-[300px] flex items-center rounded-md p-3 shadow-lg">
          <div className="p-3">
            <Workflow size={48} strokeWidth={1} className="text-emerald-500" />
          </div>
          <div>
            <div className="text-muted-foreground font-semibold">{t("dashboard.statistics.all.workflow")}</div>
            <div className="flex items-baseline">
              <div className="text-3xl text-stone-700 dark:text-stone-200">
                {statistic?.workflowTotal ? (
                  <Link to="/workflows" className="hover:underline">
                    {statistic?.workflowTotal}
                  </Link>
                ) : (
                  0
                )}
              </div>
              <div className="ml-1 text-stone-700 dark:text-stone-200">{t("dashboard.statistics.unit")}</div>
            </div>
          </div>
        </div>

        <div className="border w-full md:w-[250px] 3xl:w-[300px] flex items-center rounded-md p-3 shadow-lg">
          <div className="p-3">
            <FolderCheck size={48} strokeWidth={1} className="text-green-400" />
          </div>
          <div>
            <div className="text-muted-foreground font-semibold">{t("dashboard.statistics.enabled.workflow")}</div>
            <div className="flex items-baseline">
              <div className="text-3xl text-stone-700 dark:text-stone-200">
                {statistic?.workflowEnabled ? (
                  <Link to="/workflows?state=enabled" className="hover:underline">
                    {statistic?.workflowEnabled}
                  </Link>
                ) : (
                  0
                )}
              </div>
              <div className="ml-1 text-stone-700 dark:text-stone-200">{t("dashboard.statistics.unit")}</div>
            </div>
          </div>
        </div>
      </div>

      <div className="my-4">
        <hr />
      </div>

      <div>
        <div className="text-muted-foreground mt-5 text-sm">{t("dashboard.certificate")}</div>

        <CertificateList />
      </div>
    </div>
  );
};

export default Dashboard;
