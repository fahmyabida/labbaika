import { Divider, List, ListItem, ListItemButton, ListItemIcon, ListItemText, Toolbar } from "@mui/material"
import { Link, useLocation } from 'react-router-dom';
import InsertDriveFileIcon from '@mui/icons-material/InsertDriveFile';
import MenuBookIcon from '@mui/icons-material/MenuBook';
import CategoryIcon from '@mui/icons-material/Category';
import DashboardIcon from '@mui/icons-material/Dashboard';
import logo from '../assets/labbaika-logo.webp';

export type Page = 'dashboard' | 'cashbook' | 'category' | 'csv';

export const PAGE_TITLES: Record<Page, string> = {
  dashboard: 'Dashboard',
  cashbook: 'Cashbook',
  category: 'Categories',
  csv: 'Convert CSV',
};

export const PAGE_PATHS: Record<Page, string> = {
  dashboard: '/dashboard',
  cashbook: '/cashbook',
  category: '/category',
  csv: '/csv',
};

const items: { key: Page; label: string; icon: JSX.Element }[] = [
  { key: 'dashboard', label: PAGE_TITLES.dashboard, icon: <DashboardIcon /> },
  { key: 'cashbook', label: PAGE_TITLES.cashbook, icon: <MenuBookIcon /> },
  { key: 'category', label: PAGE_TITLES.category, icon: <CategoryIcon /> },
  { key: 'csv', label: PAGE_TITLES.csv, icon: <InsertDriveFileIcon /> },
];

interface Props {
  // Called after navigating — used to close the temporary (mobile) drawer.
  onNavigate?: () => void;
}

export const DrawerPayslip = ({ onNavigate }: Props) => {
    const location = useLocation();
    return (
      <div>
        <Toolbar sx={{ justifyContent: 'center' }}>
          <img src={logo} alt="Labbaika Catering" style={{ height: 48 }} />
        </Toolbar>
        <Divider />
        <List>
          {items.map((item) => (
            <ListItem key={item.key} disablePadding>
              <ListItemButton
                component={Link}
                to={PAGE_PATHS[item.key]}
                selected={location.pathname === PAGE_PATHS[item.key]}
                onClick={onNavigate}
              >
                <ListItemIcon>
                  {item.icon}
                </ListItemIcon>
                <ListItemText primary={item.label} />
              </ListItemButton>
            </ListItem>
          ))}
        </List>
      </div>
    )
}
