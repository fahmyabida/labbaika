import { Divider, List, ListItem, ListItemButton, ListItemIcon, ListItemText, Toolbar } from "@mui/material"
import InsertDriveFileIcon from '@mui/icons-material/InsertDriveFile';

export const DrawerPayslip = () => {
    return (
      <div>
        <Toolbar />
        <Divider />
        <List>
          {['Convert CSV'].map((text, index) => (
            <ListItem key={text} disablePadding>
              <ListItemButton>
                <ListItemIcon>
                  <InsertDriveFileIcon />
                </ListItemIcon>
                <ListItemText primary={text} />
              </ListItemButton>
            </ListItem>
          ))}
        </List>
      </div>
    )
}