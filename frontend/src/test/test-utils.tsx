import { render, type RenderOptions } from "@testing-library/react";
import { type ReactElement, type ReactNode } from "react";
import { ThemeProvider } from "next-themes";
import { I18nextProvider } from "react-i18next";
import i18n from "@/i18n";

interface WrapperProps {
  children: ReactNode;
}

function AllProviders({ children }: WrapperProps) {
  return (
    <I18nextProvider i18n={i18n}>
      <ThemeProvider attribute="class" defaultTheme="system" enableSystem>
        {children}
      </ThemeProvider>
    </I18nextProvider>
  );
}

function customRender(
  ui: ReactElement,
  options?: Omit<RenderOptions, "wrapper">
) {
  return render(ui, { wrapper: AllProviders, ...options });
}

export * from "@testing-library/react";
export { customRender as render };
export { default as userEvent } from "@testing-library/user-event";
