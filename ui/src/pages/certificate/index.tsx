import CertificateList from "@/components/certificate/CertificateList";

const Certificate = () => {
  return (
    <div className="flex flex-col space-y-5">
      <div className="text-muted-foreground">证书</div>

      <CertificateList withPagination={true} />
    </div>
  );
};

export default Certificate;
