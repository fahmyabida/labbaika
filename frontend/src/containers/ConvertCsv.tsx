import React, { useState } from 'react';
import { Button } from '@mui/material';
import { UploadPayslipCsv } from '../store/ConvertCsv';
import FileDownload from 'js-file-download'

export const ConvertCsv = () => {
    const [name, setFileName] = useState()
    const [file, setFile] = useState()

    const handleSubmit = () => {
        const result = UploadPayslipCsv(file)
        result.then((response) => {
            FileDownload(response.data, 'report.docx')
        })
    }

    const handeFileChange = (e: any) => {
        setFile(e.target.files[0])
        setFileName(e.target.value)
    } 
    return (
        <div>
            <Button variant="contained" component="label" >
                Upload File
                <input type="file" hidden onChange={handeFileChange}/>
            </Button>
            <Button onClick={handleSubmit}> Submit </Button>
            <br/>{name}
        </div>
    )
}