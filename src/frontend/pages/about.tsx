const About = () => {
  return (
    <div className="globalCss">
    <h1 className="text-5xl mt-[1%]">O nas</h1>
      <ul className="mt-5">
        <li className="text-xl pt-5">
          <b className="text-3xl">Historia fundacji:</b> Krótkie opowiadanie o
          tym, jak i dlaczego powstała fundacja, co ją napędza, jakie problemy
          próbuje rozwiązać.
        </li>
        <li className="text-xl pt-5">
          <b className="text-3xl">Misja i Wartości:</b> Podkreśl główne cele i
          wartości fundacji, takie jak pomoc społeczna, ochrona środowiska,
          edukacja, itp.
        </li>
        <li className="text-xl pt-5">
          <b className="text-3xl">Zespół:</b> Przedstawienie członków zespołu.
          Krótkie biografie z ich zdjęciami, funkcjami i być może osobistym
          cytatem, dlaczego pracują w fundacji.
        </li>
        <li className="text-xl pt-5">
          <b className="text-3xl">Współpraca:</b> Jeśli fundacja współpracuje z
          innymi organizacjami lub ma partnerów, opisz te relacje.
        </li>
      </ul>

      <br></br>
      <hr></hr>
      <br></br>

      <div>
        <h1 className="text-5xl">Kontakt</h1>

        <br></br>

        <h2>Adres:</h2>
        <p>
          <b>Ul.partosumadre23</b>
          <br />
          <b>Pon - Pt :</b> 10:00 - 18:00
          <br />
          <b>Sobota :</b> 09:00 - 14:00
        </p>
        <br></br>
        <h2>Kontakt:</h2>
        <p>
          <b>Tel:</b> +48 000 000 000
          <br />
          <b>E-mail:</b> dsadosadasd@gmail.com
        </p>
<hr></hr>
        <iframe
          title="b"
          src="https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d2557.0223748554645!2d19.409808976539107!3d50.14201327153497!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x4716eeb11871fa5d%3A0xe3333daf7f955250!2sCentrum%20Handlowe%20MAX!5e0!3m2!1spl!2spl!4v1697022456514!5m2!1spl!2spl"
          allowFullScreen={false}
          loading="lazy"
          referrerPolicy="no-referrer-when-downgrade"
          className="w-full h-[80vh]  py-5"
        ></iframe>
      </div>
    </div>
  );
};
export default About;
