// MisViáticos Landing Page - Type Definitions

export interface NavigationItem {
  label: string;
  href: string;
  hasDropdown?: boolean;
  dropdownItems?: NavigationItem[];
}

export interface Feature {
  id: string;
  title: string;
  description: string;
  icon: React.ComponentType<{ className?: string }>;
  benefits: string[];
}

export interface CTAProps {
  title: string;
  subtitle: string;
  primaryButton: {
    text: string;
    action: () => void;
  };
  secondaryButton?: {
    text: string;
    action: () => void;
  };
}

export interface TestimonialProps {
  name: string;
  company: string;
  role: string;
  content: string;
  avatar?: string;
  rating: number;
}

export interface HeroProps {
  title: string;
  subtitle: string;
  ctaButtons: {
    primary: string;
    secondary: string;
  };
}

export interface FooterLink {
  label: string;
  href: string;
}

export interface FooterSection {
  title: string;
  links: FooterLink[];
}
