import { Menu } from "@suid/icons-material";
import {
  AppBar,
  Box,
  Button,
  IconButton,
  Toolbar,
  Typography,
} from "@suid/material";
import type { JSX } from "solid-js";
import routes from "../utils/routes";

export default function Layout({ children }: { children: JSX.Element }) {
  return (
    <Box sx={{ flexGrow: 1 }}>
      <AppBar position="static">
        <Toolbar>
          <IconButton
            size="large"
            edge="start"
            color="inherit"
            aria-label="menu"
            sx={{ mr: 2 }}
          >
            <Menu />
          </IconButton>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            Karaberus
          </Typography>
          <Button color="inherit" href={routes.OIDC_LOGIN}>
            Login
          </Button>
        </Toolbar>
      </AppBar>
      {children}
    </Box>
  );
}
