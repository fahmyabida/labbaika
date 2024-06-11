import axios from "axios";
import { PAYSLIP_SERVICE_URL } from ".";

export async function UploadPayslipCsv(file: any)  {
    var formData = new FormData();
    formData.append("file", file);
    var result : any
    try {
        result = await axios.post(`${PAYSLIP_SERVICE_URL}/api/v1/payslips/convert`, formData, {
            responseType: 'blob',
            headers: {
                'Content-Type': 'multipart/form-data',
                'Accept': '*/*'
            }
        })
        console.log("result UploadPayslipCsv", result)
    } catch (err) {
        console.log("error while UploadPayslipCsv", err)
    } finally {
        return result
    }
}