import Header from './sections/Header'
import Hero from './sections/Hero'
import Features from './sections/Features'
import CTA from './sections/CTA'
import Footer from './sections/Footer'

export default function LandingPage() {
  return (
    <div className="min-h-screen bg-white">
      <Header />
      <Hero />
      <Features />
      <CTA />
      <Footer />
    </div>
  )
}
