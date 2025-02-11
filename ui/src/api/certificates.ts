import { ClientResponseError } from "pocketbase";

import { type CertificateFormatType } from "@/domain/certificate";
import { getPocketBase } from "@/repository/_pocketbase";

type ArchiveRespData = {
  fileBytes: string;
};

export const archive = async (certificateId: string, format?: CertificateFormatType) => {
  const pb = getPocketBase();

  const resp = await pb.send<BaseResponse<ArchiveRespData>>(`/api/certificates/${encodeURIComponent(certificateId)}/archive`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: {
      format: format,
    },
  });

  if (resp.code != 0) {
    throw new ClientResponseError({ status: resp.code, response: resp, data: {} });
  }

  return resp;
};

type ValidateCertificateResp = {
  isValid: boolean;
  domains: string;
};

export const validateCertificate = async (certificate: string) => {
  const pb = getPocketBase();
  const resp = await pb.send<BaseResponse<ValidateCertificateResp>>(`/api/certificates/validate/certificate`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: {
      certificate: certificate,
    },
  });

  if (resp.code != 0) {
    throw new ClientResponseError({ status: resp.code, response: resp, data: {} });
  }

  return resp;
};

type ValidatePrivateKeyResp = {
  isValid: boolean;
};

export const validatePrivateKey = async (privateKey: string) => {
  const pb = getPocketBase();
  const resp = await pb.send<BaseResponse<ValidatePrivateKeyResp>>(`/api/certificates/validate/private-key`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: {
      privateKey: privateKey,
    },
  });

  if (resp.code != 0) {
    throw new ClientResponseError({ status: resp.code, response: resp, data: {} });
  }

  return resp;
};
