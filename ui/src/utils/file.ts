import JSZip from "jszip";

export function readFileContent(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();

    reader.onload = () => {
      if (reader.result) {
        resolve(reader.result.toString());
      } else {
        reject("No content found");
      }
    };

    reader.onerror = () => reject(reader.error);

    reader.readAsText(file, "utf-8");
  });
}

export const saveFiles2Zip = async (zipName: string, files: { name: string; content: string }[]) => {
  const zip = new JSZip();

  files.forEach((file) => {
    zip.file(file.name, file.content);
  });

  const content = await zip.generateAsync({ type: "blob" });

  // Save the zip file to the local system
  const link = document.createElement("a");
  link.href = URL.createObjectURL(content);
  link.download = zipName;
  link.click();
};
