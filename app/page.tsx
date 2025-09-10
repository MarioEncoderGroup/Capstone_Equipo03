import Header from './_landingpage/components/Header'
import Hero from './_landingpage/components/Hero'
import Features from './_landingpage/components/Features'
import CTA from './_landingpage/components/CTA'
import Footer from './_landingpage/components/Footer'

export default function Page() {
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
