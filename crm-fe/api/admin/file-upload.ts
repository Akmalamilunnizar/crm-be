import type { CreateFileUploadRequest } from "@/types/requests/file-upload";

export const uploadFileAdminApi = () => {
  const api = useApiHost();
  return {
    createUploadFile: async (
      data: CreateFileUploadRequest
    ) => {
      
  const formData = new FormData();
  formData.append("name", data.name);
  formData.append("path", data.path);
  formData.append("file", data.file); // pa
      const response = await fetch(`${api}/api/file-upload`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${useCookie("token").value}`,
        },
        body: formData,
      });
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Login failed');
      }
      return response.json();
    },
  };
};