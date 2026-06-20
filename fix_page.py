content = """\"use client\";

import Link from \"next/link\";
import Image from \"next/image\";
import { useEffect, useRef, useState } from \"react\";

const allRepos = [
  {
    name: \"whisper\",
    desc: \"Hell\x27s Gate / Halo\x27s Gate for Linux. Indirect syscalls via runtime libc ELF parsing.\",
    cve: null,
    lang: \"Go\",
    url: \"https://github.com/Koshmare-Blossom/whisper\",
    stars: 0,
  },
  {
    name: \"eclipse\",
    desc: \"Linux Sleep Obfuscation.\",
    cve: null,
    lang: \"Go\",
    url: \"https://github.com/Koshmare-Blossom/eclipse\",
    stars: 0,
  },
  {
    name: \"Dear-Linux-With-Love\",
    desc: \"Linux kernel exploits, rewritten with love.\",
    cve: null,
    lang: \"Go\",
    url: \"https://github.com/Koshmare-Blossom/Dear-Linux-With-Love\",
    stars: 0,
  },
  {
    name: \"PinTheft-go\",
    desc: \"A Go implementation of PinTheft (CVE-2026-43494)\",
    cve: \"CVE-2026-43494\",
    lang: \"Go\",
    url: \"https://github.com/Koshmare-Blossom/PinTheft-go\",
    stars: 0,
  },
  {
    name: \"PinTheft-asm\",
    desc: \"A x86_64 ASM implementation of PinTheft (CVE-2026-43494)\",
    cve: \"CVE-2026-43494\",
    lang: \"ASM\",
    url: \"https://github.com/Koshmare-Blossom/PinTheft-asm\",
    stars: 0,
  },
  {
    name: \"DirtyFrag-go\",
    desc: \"A Go implementation of dirtyfrag (CVE-2026-43284 / CVE-2026-43500)\",
    cve: \"CVE-2026-43284 / CVE-2026-43500\",
    lang: \"Go\",
    url: \"https://github.com/Koshmare-Blossom/DirtyFrag-go\",
    stars: 0,
  },
  {
    name: \"DirtyDecrypt-go\",
    desc: \"A Go implementation of dirtydecrypt (CVE-2026-31635)\",
    cve: \"CVE-2026-31635\",
    lang: \"Go\",
    url: \"https://github.com/Koshmare-Blossom/DirtyDecrypt-go\",
    stars: 0,
  },
  {
    name: \"Fragnesia-go\",
    desc: \"A Go implementation of fragnesia (CVE-2026-46300)\",
    cve: \"CVE-2026-46300\",
    lang: \"Go\",
    url: \"https://github.com/Koshmare-Blossom/Fragnesia-go\",
    stars: 0,
  },
  {
    name: \"CIFSwitch-go\",
    desc: \"A Go implementation of CIFSwitch (CVE-2026-46243)\",
    cve: \"CVE-2026-46243\",
    lang: \"Go\",
    url: \"https://github.com/Koshmare-Blossom/CIFSwitch-go\",
    stars: 0,
  },
  {
    name: \"Copyfail-sh\",
    desc: \"A Bash implementation of copyfail (CVE-2026-31431)\",
    cve: \"CVE-2026-31431\",
    lang: \"Shell\",
    url: \"https://github.com/Koshmare-Blossom/Copyfail-sh\",
    stars: 0,
  },
];

const langColor: Record<string, string> = {
  Go: \"#00acd7\",
  ASM: \"#a78bfa\",
  Shell: \"#e879f9\",
};

export default function Home() {
  const heroRef = useRef<HTMLDivElement>(null);
  const [featured, setFeatured] = useState<any[]>([]);

  useEffect(() => {
    const el = heroRef.current;
    if (!el) return;
    el.style.opacity = \"0\";
    el.style.transform = \"translateY(16px)\";
    requestAnimationFrame(() => {
      el.style.transition = \"opacity 0.6s ease, transform 0.6s ease\";
      el.style.opacity = \"1\";
      el.style.transform = \"translateY(0)\";
    });

    // Pick 4 random repos
    const shuffled = [...allRepos].sort(() => 0.5 - Math.random());
    const selected = shuffled.slice(0, 4);
    setFeatured(selected);

    const fetchStars = async () => {
      try {
        const updated = await Promise.all(
          selected.map(async (repo) => {
            const res = await fetch(\"https://api.github.com/repos/Koshmare-Blossom/\" + repo.name);
            if (res.ok) {
              const data = await res.json();
              return { ...repo, stars: data.stargazers_count };
            }
            return repo;
          })
        );
        setFeatured(updated);
      } catch (err) {
        console.error(err);
      }
    };
    fetchStars();
  }, []);

  return (
    <div className=\"max-w-6xl mx-auto px-6 py-20\">

      {/* Hero */}
      <section ref={heroRef} className=\"mb-24 flex items-start gap-8\">
        <div className=\"flex-1\">
          <p className=\"font-mono text-sm text-[#64748b] mb-4\">
            linux kernel research
          </p>
          <h1 className=\"text-4xl font-semibold text-[#e2e8f0] mb-4 glow-pink\">
            Koshmare-Blossom
          </h1>
          <p className=\"text-[#94a3b8] text-lg max-w-xl leading-relaxed\">
            I spend my time reimplementing CVEs targeting the Linux kernel -
            in different languages, to understand exactly how they work
            and what they expose.
          </p>
          <div className=\"flex items-center gap-4 mt-8\">
            <Link
              href=\"/research\"
              className=\"font-mono text-sm px-4 py-2 border border-[#e879f933] text-[#e879f9] hover:bg-[#e879f910] rounded transition-colors\"
            >
              view research
            </Link>
            <Link
              href=\"/blog\"
              className=\"font-mono text-sm text-[#64748b] hover:text-[#e2e8f0] transition-colors\"
            >
              read the blog →
            </Link>
          </div>
        </div>
        <div className=\"shrink-0 hidden sm:block\">
          <Image
            src=\"/avatar.jpg\"
            alt=\"Koshmare-Blossom\"
            width={160}
            height={160}
            className=\"rounded-full border border-[#1e1e30]\"
            priority
          />
        </div>
      </section>

      {/* Featured repos */}
      <section>
        <div className=\"flex items-center justify-between mb-6\">
          <h2 className=\"font-mono text-xs text-[#64748b] uppercase tracking-widest\">
            featured
          </h2>
          <Link
            href=\"/research\"
            className=\"font-mono text-xs text-[#64748b] hover:text-[#e879f9] transition-colors\"
          >
            view all →
          </Link>
        </div>
        <div className=\"grid grid-cols-1 sm:grid-cols-2 gap-3\">
          {featured.map((repo) => (
            <a
              key={repo.name}
              href={repo.url}
              target=\"_blank\"
              rel=\"noopener noreferrer\"
              className=\"group block p-4 border border-[#1e1e30] rounded-lg bg-[#0f0f1a] hover:border-[#e879f933] transition-all\"
            >
              <div className=\"flex items-start justify-between mb-2\">
                <div className=\"flex items-center gap-2\">
                  <span className=\"font-mono text-sm text-[#e2e8f0] group-hover:text-[#f0abfc] transition-colors\">
                    {repo.name}
                  </span>
                  {repo.stars > 0 && (
                    <span className=\"font-mono text-[10px] text-[#64748b]\">
                       ★ {repo.stars}
                    </span>
                  )}
                </div>
                <span
                  className=\"font-mono text-xs px-1.5 py-0.5 rounded\"
                  style={{
                    color: langColor[repo.lang] ?? \"#94a3b8\",
                    background: (langColor[repo.lang] ?? \"#94a3b8\") + \"15\",
                  }}
                >
                  {repo.lang}
                </span>
              </div>
              <p className=\"text-[#64748b] text-xs leading-relaxed mb-2\">
                {repo.desc}
              </p>
              {repo.cve && (
                <span className=\"font-mono text-xs text-[#a78bfa]\">
                  {repo.cve}
                </span>
              )}
            </a>
          ))}
        </div>
      </section>

    </div>
  );
}
\"\"\"
with open(\"/home/bbuddha/github/Koshmare-Blossom.github.io/app/page.tsx\", \"w\") as f:
    f.write(content.replace(\"\\\\x27\", \"\x27\"))
