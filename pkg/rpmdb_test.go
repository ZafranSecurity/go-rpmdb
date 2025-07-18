package rpmdb

import (
	"encoding/hex"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/glebarez/go-sqlite"
)

var packageTests = []struct {
	name    string
	file    string // Test input file
	pkgList []*PackageInfo
}{
	{
		name:    "CentOS5 plain",
		file:    "testdata/centos5-plain/Packages",
		pkgList: CentOS5Plain(),
	},
	{
		name:    "CentOS6 Plain",
		file:    "testdata/centos6-plain/Packages",
		pkgList: CentOS6Plain(),
	},
	{
		name:    "CentOS6 with Development tools",
		file:    "testdata/centos6-devtools/Packages",
		pkgList: CentOS6DevTools(),
	},
	{
		name:    "CentOS6 with many packages",
		file:    "testdata/centos6-many/Packages",
		pkgList: CentOS6Many(),
	},
	{
		name:    "CentOS7 Plain",
		file:    "testdata/centos7-plain/Packages",
		pkgList: CentOS7Plain(),
	},
	{
		name:    "CentOS7 with Development tools",
		file:    "testdata/centos7-devtools/Packages",
		pkgList: CentOS7DevTools(),
	},
	{
		name:    "CentOS7 with many packages",
		file:    "testdata/centos7-many/Packages",
		pkgList: CentOS7Many(),
	},
	{
		name:    "CentOS7 with Python 3.5",
		file:    "testdata/centos7-python35/Packages",
		pkgList: CentOS7Python35(),
	},
	{
		name:    "CentOS7 with httpd 2.4",
		file:    "testdata/centos7-httpd24/Packages",
		pkgList: CentOS7Httpd24(),
	},
	{
		name:    "CentOS8 with modules",
		file:    "testdata/centos8-modularitylabel/Packages",
		pkgList: CentOS8Modularitylabel(),
	},
	{
		name:    "RHEL UBI8 from s390x",
		file:    "testdata/ubi8-s390x/Packages",
		pkgList: UBI8s390x(),
	},
	{
		name:    "SLE15 with NDB style rpm database",
		file:    "testdata/sle15-bci/Packages.db",
		pkgList: SLE15WithNDB(),
	},
	{
		name:    "Fedora35 with SQLite3 style rpm database",
		file:    "testdata/fedora35/rpmdb.sqlite",
		pkgList: Fedora35WithSQLite3(),
	},
	{
		name:    "Fedora35 plus MongoDB with SQLite3 style rpm database",
		file:    "testdata/fedora35-plus-mongo/rpmdb.sqlite",
		pkgList: Fedora35PlusMongoDBWithSQLite3(),
	},
	{
		name:    "Rocky9 with SQLite3 style rpm database (newer signature format)",
		file:    "testdata/rockylinux-9/rpmdb.sqlite",
		pkgList: Rockylinux9WithSQLite3(),
	},
}

func TestPackageList(t *testing.T) {
	for _, tt := range packageTests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := Open(tt.file)
			require.NoError(t, err)

			got, err := db.ListPackages()
			require.NoError(t, err)

			// They are tested in another function.
			for _, g := range got {
				g.PGP = ""
				g.RSAHeader = ""
				g.DigestAlgorithm = 0
				g.InstallTime = 0
				g.BaseNames = nil
				g.DirIndexes = nil
				g.DirNames = nil
				g.FileSizes = nil
				g.FileDigests = nil
				g.FileModes = nil
				g.FileFlags = nil
				g.UserNames = nil
				g.GroupNames = nil
				g.Provides = nil
				g.Requires = nil
			}

			for i, p := range tt.pkgList {
				assert.Equal(t, p, got[i])
			}
		})
	}
}

func BenchmarkRpmDB_Package(b *testing.B) {
	for _, tt := range packageTests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				db, err := Open(tt.file)
				if err != nil {
					b.Fatal(err)
				}
				_, err = db.ListPackages()
				if err != nil {
					b.Fatal(err)
				}
			}
			b.ReportAllocs()
		})
	}
}

func Test_parseRSA(t *testing.T) {
	tests := []struct {
		name    string
		ie      IndexEntry
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "older RSA header",
			ie: IndexEntry{
				Data: func() []byte {
					h := "89021503050058d3e39b0946fca2c105b9de0102b12a1000a2b3d347b51142e83b2de5e03ba9096f6330b72c140e46200d662b01c78534d14fab2ad4f07325119386830dd590219f27a22e420680283c500c40e6fba95404884b0a0abca8f198030ddc03653b7db2883b8230687e9e73d43eb5a24dbabfa48bbb3d1151ed264744e5e8ca169b0c4673a1440a9b99e53e693c9722f6423833cd7795e3044227fb922e21b7c007f03e923fae3f04d1ac2e8581e68c6790115b6dccfc02c8cb41681ed84785df086d6e26008c257d088a524ba2e7a7a5f41ad26b106c67b87fe48118b69662db612c23d2140059286f1ba7764627def6867ad0e11fe3a01fb1422dabe6f5cdf4cd876dc4fadfd2364bc3ba3758db94aaf3b82368cba65cf762287f713eb7ddc773acf93b083c739577a7eaf1f99e7dcbb8db1da050490e9fb67c838448db060a9e619d318c96f03e4363808d84ce29e8c102c290cc2bfab5746f3d9ddc9eb8b428f3ad2678abb2d46e846ddca7fc41322d76a97be6d416b4750f23320ec725e082be4496483b4cd3a3d2c515b3c8a6e27541139d809245140303877b84842ed2dd0454a78b2dfb7d6213784697077a8167942ebda5995a28d8256957e33e301706c35944ae05c7a54a4dd89be654d26cefa5cf0f616bbeaf317138371b09c5bbd5531f716020e553354ce5dbce3d9bb72f21e1857408dfd5a35250ff40f61ae1e25409ae9d21db76b8878341f4762a22be2189"
					b, err := hex.DecodeString(h)
					if err != nil {
						t.Fatalf("failed to decode hex string: %v", err)
					}
					return b
				}(),
			},
			want: "RSA/SHA1, Thu Mar 23 15:02:51 2017, Key ID 0946fca2c105b9de",
		},
		{
			name: "newer RSA header",
			ie: IndexEntry{
				Data: func() []byte {
					h := "89024a04000108003416210421cb256ae16fc54c6e652949702d426d350d275d050262804369161c72656c656e6740726f636b796c696e75782e6f7267000a0910702d426d350d275dc8910ffd14f0f80297481fea648e7ba5a74bce10c5faccc2bbe588caece04be34d304a6a445538afc97a7033d43c983d27cc8f5ee515b2dd92f3e03354c413e55372a4d19386eb0f2354f9a26ee5fc2e56dfda49555e4a58b49279b70cd2036b04f28125f85942f640f2984e29e079f26bf6f76831d83d95983aa084a3e7b6327be2e23d0d799c4b4d1cfb36147ddfb782bf9df7b331d97f4f46b38f968b6130d87b0ef6bb0d424390fe34e38092babed37440569a93f55f50a2bdb58be0259f35badf7e728bd49824ed47f69fa53b6e26736bde4d8358d959b090e88054c3e179745dc7377e41b54b4e10223f4859e88162c7c5ec64b78d36cf8a914c1c2deb8c4f19a70d406e70756a89195d6aee488a9b40b9dbb76b2c38e528eb88d08ec35774a48ed9ce4e0dfac45cb7613ad5921f54c61d3aae5d7b3ab0e2e6ff867ac8f395b37af78b5c01022a4a4e62f7a99425fccb7439880cd6b393a3050b2e9512693bc36f6fe9de2921dda59710a1508965065244cf9f0f8cfc5bd554777f1a84d2249339234d62f2441249f617ad7df4fb01367a91d3a880e86fdb84bc6d03a127b44a28c6ceadef89e438db9640aa59b8a3f460b07272511f8187a5f3b163c8fd1caa61667401bce2ccdb1c176c46be10ef8033903132cca5889fa3661b2fba590c41fa1c104c08426677bdbf745a52ccd28f581960cf9d7e4ede3b9584aacb2f20ef93"
					b, err := hex.DecodeString(h)
					if err != nil {
						t.Fatalf("failed to decode hex string: %v", err)
					}
					return b
				}(),
			},
			want: "RSA/SHA256, Sun May 15 00:03:53 2022, Key ID 702d426d350d275d",
		},
		{
			// $ curl -O https://download.postgresql.org/pub/repos/yum/14/redhat/rhel-9-x86_64/postgresql14-server-14.10-1PGDG.rhel9.x86_64.rpm
			// $ rpm -ivh --nodeps --force postgresql14-server-14.10-1PGDG.rhel9.x86_64.rpm
			// $ rpm -q --qf '%{NAME}-%{VERSION}-%{RELEASE} %{RSAHEADER:pgpsig}\n' postgresql14-server-14.10-1PGDG.rhel9.x86_64
			//   postgresql14-server-14.10-1PGDG.rhel9 RSA/SHA256, Tue Jan  2 16:45:56 2024, Key ID 40bca2b408b40d20
			name: "example from rocky9 postgresql server",
			ie: IndexEntry{
				Data: func() []byte {
					h := "8901b304000108001d162104d4bf08ae67a0b4c7a1dbccd240bca2b408b40d20050265943dc4000a091040bca2b408b40d203b270bff71678ffeb190833a19a82112f59eee64cba186ab454d4526e0b3c8797e723f6916daff1b1f18cbf53c0da5d398a3a42065e79e5ca939f721652f38400dd4cac1107a902b1dae880649437ad0242444f3f07115172cae0a207b7cf8340af2f4a94976325f1dc165d5c2a564be322c4e130adb6217e7138b689f08898c407b223aa1ff8f8d592f31eba2256c02fae70ce4022d688a487972646b8bf1b518b5d6549c1e60fd812134422d9fdb41cf799f5eab80e48b4ab7cff84362dc867ed1af1416dd78e92bcc59217de7064b9a015d94a5097788689b9b6fbdeea679cfe4a6947f73dc3a6c810f2cb999d279b01564422d1500fc1bd8bd1eefa2d60660127ffef24067354660f93c0faf81f4edd599dd7e4b77fe4bff6c7a0ea83530c817c38d1f2364175883c6ef7b6dec86ad282bdd5138b8597567db96810c4ed6454a4ab1d98f0425dcd8892a5d46ed9289cb3ae3e1f1e2663d3e8188e873428f6cf7163563ed3860edc4fee81522389508847e692e2d13310eb4b40f7fdd7eb364a0b2dc"
					b, err := hex.DecodeString(h)
					if err != nil {
						t.Fatalf("failed to decode hex string: %v", err)
					}
					return b
				}(),
			},
			want: "RSA/SHA256, Tue Jan  2 16:45:56 2024, Key ID 40bca2b408b40d20",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr == nil {
				tt.wantErr = assert.NoError
			}
			got, err := parsePGP(tt.ie)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRpmDB_Package(t *testing.T) {
	tests := []struct {
		name                   string
		pkgName                string
		file                   string // Test input file
		want                   *PackageInfo
		wantInstalledFiles     []FileInfo
		wantInstalledFileNames []string
		wantErr                string
	}{
		{
			name:    "centos5 python",
			pkgName: "python",
			file:    "testdata/centos5-plain/Packages",
			want: &PackageInfo{
				Name:        "python",
				Version:     "2.4.3",
				Release:     "56.el5",
				Arch:        "x86_64",
				Size:        74377,
				SourceRpm:   "python-2.4.3-56.el5.src.rpm",
				License:     "PSF - see LICENSE",
				Vendor:      "CentOS",
				Summary:     "An interpreted, interactive, object-oriented programming language.",
				SigMD5:      "ebfb56be33b146ef39180a090e581258",
				PGP:         "",
				InstallTime: 1459411575,
				Provides: []string{
					"Distutils",
					"python(abi)",
					"python-abi",
					"python-x86_64",
					"python2",
					"python",
				},
				Requires: []string{
					"/usr/bin/env",
					"libc.so.6()(64bit)",
					"libc.so.6(GLIBC_2.2.5)(64bit)",
					"libdl.so.2()(64bit)",
					"libm.so.6()(64bit)",
					"libpthread.so.0()(64bit)",
					"libpython2.4.so.1.0()(64bit)",
					"libutil.so.1()(64bit)",
					"python-libs-x86_64",
					"rpmlib(CompressedFileNames)",
					"rpmlib(PartialHardlinkSets)",
					"rpmlib(PayloadFilesHavePrefix)",
					"rpmlib(VersionedDependencies)",
					"rtld(GNU_HASH)",
				},
			},
			wantInstalledFiles:     CentOS5PythonInstalledFiles,
			wantInstalledFileNames: CentOS5PythonInstalledFileNames,
		},
		{
			name:    "centos6 glibc",
			pkgName: "glibc",
			file:    "testdata/centos6-plain/Packages",
			want: &PackageInfo{
				Name:            "glibc",
				Version:         "2.12",
				Release:         "1.212.el6",
				Arch:            "x86_64",
				Size:            13117447,
				SourceRpm:       "glibc-2.12-1.212.el6.src.rpm",
				License:         "LGPLv2+ and LGPLv2+ with exceptions and GPLv2+",
				Vendor:          "CentOS",
				Summary:         "The GNU libc libraries",
				SigMD5:          "89e843d7979a50a26e2ea1924ef3e213",
				DigestAlgorithm: PGPHASHALGO_SHA256,
				PGP:             "RSA/SHA1, Wed Jun 20 11:36:27 2018, Key ID 0946fca2c105b9de",
				RSAHeader:       "RSA/SHA1, Wed Jun 20 11:36:27 2018, Key ID 0946fca2c105b9de",
				InstallTime:     1538857091,
				Provides: []string{
					"ANSI_X3.110.so()(64bit)",
					"ARMSCII-8.so()(64bit)",
					"ASMO_449.so()(64bit)",
					"BIG5.so()(64bit)",
					"BIG5HKSCS.so()(64bit)",
					"BRF.so()(64bit)",
					"CP10007.so()(64bit)",
					"CP1125.so()(64bit)",
					"CP1250.so()(64bit)",
					"CP1251.so()(64bit)",
					"CP1252.so()(64bit)",
					"CP1253.so()(64bit)",
					"CP1254.so()(64bit)",
					"CP1255.so()(64bit)",
					"CP1256.so()(64bit)",
					"CP1257.so()(64bit)",
					"CP1258.so()(64bit)",
					"CP737.so()(64bit)",
					"CP775.so()(64bit)",
					"CP932.so()(64bit)",
					"CSN_369103.so()(64bit)",
					"CWI.so()(64bit)",
					"DEC-MCS.so()(64bit)",
					"EBCDIC-AT-DE-A.so()(64bit)",
					"EBCDIC-AT-DE.so()(64bit)",
					"EBCDIC-CA-FR.so()(64bit)",
					"EBCDIC-DK-NO-A.so()(64bit)",
					"EBCDIC-DK-NO.so()(64bit)",
					"EBCDIC-ES-A.so()(64bit)",
					"EBCDIC-ES-S.so()(64bit)",
					"EBCDIC-ES.so()(64bit)",
					"EBCDIC-FI-SE-A.so()(64bit)",
					"EBCDIC-FI-SE.so()(64bit)",
					"EBCDIC-FR.so()(64bit)",
					"EBCDIC-IS-FRISS.so()(64bit)",
					"EBCDIC-IT.so()(64bit)",
					"EBCDIC-PT.so()(64bit)",
					"EBCDIC-UK.so()(64bit)",
					"EBCDIC-US.so()(64bit)",
					"ECMA-CYRILLIC.so()(64bit)",
					"EUC-CN.so()(64bit)",
					"EUC-JISX0213.so()(64bit)",
					"EUC-JP-MS.so()(64bit)",
					"EUC-JP.so()(64bit)",
					"EUC-KR.so()(64bit)",
					"EUC-TW.so()(64bit)",
					"GB18030.so()(64bit)",
					"GBBIG5.so()(64bit)",
					"GBGBK.so()(64bit)",
					"GBK.so()(64bit)",
					"GEORGIAN-ACADEMY.so()(64bit)",
					"GEORGIAN-PS.so()(64bit)",
					"GOST_19768-74.so()(64bit)",
					"GREEK-CCITT.so()(64bit)",
					"GREEK7-OLD.so()(64bit)",
					"GREEK7.so()(64bit)",
					"HP-GREEK8.so()(64bit)",
					"HP-ROMAN8.so()(64bit)",
					"HP-ROMAN9.so()(64bit)",
					"HP-THAI8.so()(64bit)",
					"HP-TURKISH8.so()(64bit)",
					"IBM037.so()(64bit)",
					"IBM038.so()(64bit)",
					"IBM1004.so()(64bit)",
					"IBM1008.so()(64bit)",
					"IBM1008_420.so()(64bit)",
					"IBM1025.so()(64bit)",
					"IBM1026.so()(64bit)",
					"IBM1046.so()(64bit)",
					"IBM1047.so()(64bit)",
					"IBM1097.so()(64bit)",
					"IBM1112.so()(64bit)",
					"IBM1122.so()(64bit)",
					"IBM1123.so()(64bit)",
					"IBM1124.so()(64bit)",
					"IBM1129.so()(64bit)",
					"IBM1130.so()(64bit)",
					"IBM1132.so()(64bit)",
					"IBM1133.so()(64bit)",
					"IBM1137.so()(64bit)",
					"IBM1140.so()(64bit)",
					"IBM1141.so()(64bit)",
					"IBM1142.so()(64bit)",
					"IBM1143.so()(64bit)",
					"IBM1144.so()(64bit)",
					"IBM1145.so()(64bit)",
					"IBM1146.so()(64bit)",
					"IBM1147.so()(64bit)",
					"IBM1148.so()(64bit)",
					"IBM1149.so()(64bit)",
					"IBM1153.so()(64bit)",
					"IBM1154.so()(64bit)",
					"IBM1155.so()(64bit)",
					"IBM1156.so()(64bit)",
					"IBM1157.so()(64bit)",
					"IBM1158.so()(64bit)",
					"IBM1160.so()(64bit)",
					"IBM1161.so()(64bit)",
					"IBM1162.so()(64bit)",
					"IBM1163.so()(64bit)",
					"IBM1164.so()(64bit)",
					"IBM1166.so()(64bit)",
					"IBM1167.so()(64bit)",
					"IBM12712.so()(64bit)",
					"IBM1364.so()(64bit)",
					"IBM1371.so()(64bit)",
					"IBM1388.so()(64bit)",
					"IBM1390.so()(64bit)",
					"IBM1399.so()(64bit)",
					"IBM16804.so()(64bit)",
					"IBM256.so()(64bit)",
					"IBM273.so()(64bit)",
					"IBM274.so()(64bit)",
					"IBM275.so()(64bit)",
					"IBM277.so()(64bit)",
					"IBM278.so()(64bit)",
					"IBM280.so()(64bit)",
					"IBM281.so()(64bit)",
					"IBM284.so()(64bit)",
					"IBM285.so()(64bit)",
					"IBM290.so()(64bit)",
					"IBM297.so()(64bit)",
					"IBM420.so()(64bit)",
					"IBM423.so()(64bit)",
					"IBM424.so()(64bit)",
					"IBM437.so()(64bit)",
					"IBM4517.so()(64bit)",
					"IBM4899.so()(64bit)",
					"IBM4909.so()(64bit)",
					"IBM4971.so()(64bit)",
					"IBM500.so()(64bit)",
					"IBM5347.so()(64bit)",
					"IBM803.so()(64bit)",
					"IBM850.so()(64bit)",
					"IBM851.so()(64bit)",
					"IBM852.so()(64bit)",
					"IBM855.so()(64bit)",
					"IBM856.so()(64bit)",
					"IBM857.so()(64bit)",
					"IBM860.so()(64bit)",
					"IBM861.so()(64bit)",
					"IBM862.so()(64bit)",
					"IBM863.so()(64bit)",
					"IBM864.so()(64bit)",
					"IBM865.so()(64bit)",
					"IBM866.so()(64bit)",
					"IBM866NAV.so()(64bit)",
					"IBM868.so()(64bit)",
					"IBM869.so()(64bit)",
					"IBM870.so()(64bit)",
					"IBM871.so()(64bit)",
					"IBM874.so()(64bit)",
					"IBM875.so()(64bit)",
					"IBM880.so()(64bit)",
					"IBM891.so()(64bit)",
					"IBM901.so()(64bit)",
					"IBM902.so()(64bit)",
					"IBM903.so()(64bit)",
					"IBM9030.so()(64bit)",
					"IBM904.so()(64bit)",
					"IBM905.so()(64bit)",
					"IBM9066.so()(64bit)",
					"IBM918.so()(64bit)",
					"IBM921.so()(64bit)",
					"IBM922.so()(64bit)",
					"IBM930.so()(64bit)",
					"IBM932.so()(64bit)",
					"IBM933.so()(64bit)",
					"IBM935.so()(64bit)",
					"IBM937.so()(64bit)",
					"IBM939.so()(64bit)",
					"IBM943.so()(64bit)",
					"IBM9448.so()(64bit)",
					"IEC_P27-1.so()(64bit)",
					"INIS-8.so()(64bit)",
					"INIS-CYRILLIC.so()(64bit)",
					"INIS.so()(64bit)",
					"ISIRI-3342.so()(64bit)",
					"ISO-2022-CN-EXT.so()(64bit)",
					"ISO-2022-CN.so()(64bit)",
					"ISO-2022-JP-3.so()(64bit)",
					"ISO-2022-JP.so()(64bit)",
					"ISO-2022-KR.so()(64bit)",
					"ISO-IR-197.so()(64bit)",
					"ISO-IR-209.so()(64bit)",
					"ISO646.so()(64bit)",
					"ISO8859-1.so()(64bit)",
					"ISO8859-10.so()(64bit)",
					"ISO8859-11.so()(64bit)",
					"ISO8859-13.so()(64bit)",
					"ISO8859-14.so()(64bit)",
					"ISO8859-15.so()(64bit)",
					"ISO8859-16.so()(64bit)",
					"ISO8859-2.so()(64bit)",
					"ISO8859-3.so()(64bit)",
					"ISO8859-4.so()(64bit)",
					"ISO8859-5.so()(64bit)",
					"ISO8859-6.so()(64bit)",
					"ISO8859-7.so()(64bit)",
					"ISO8859-8.so()(64bit)",
					"ISO8859-9.so()(64bit)",
					"ISO8859-9E.so()(64bit)",
					"ISO_10367-BOX.so()(64bit)",
					"ISO_11548-1.so()(64bit)",
					"ISO_2033.so()(64bit)",
					"ISO_5427-EXT.so()(64bit)",
					"ISO_5427.so()(64bit)",
					"ISO_5428.so()(64bit)",
					"ISO_6937-2.so()(64bit)",
					"ISO_6937.so()(64bit)",
					"JOHAB.so()(64bit)",
					"KOI-8.so()(64bit)",
					"KOI8-R.so()(64bit)",
					"KOI8-RU.so()(64bit)",
					"KOI8-T.so()(64bit)",
					"KOI8-U.so()(64bit)",
					"LATIN-GREEK-1.so()(64bit)",
					"LATIN-GREEK.so()(64bit)",
					"MAC-CENTRALEUROPE.so()(64bit)",
					"MAC-IS.so()(64bit)",
					"MAC-SAMI.so()(64bit)",
					"MAC-UK.so()(64bit)",
					"MACINTOSH.so()(64bit)",
					"MIK.so()(64bit)",
					"NATS-DANO.so()(64bit)",
					"NATS-SEFI.so()(64bit)",
					"PT154.so()(64bit)",
					"RK1048.so()(64bit)",
					"SAMI-WS2.so()(64bit)",
					"SHIFT_JISX0213.so()(64bit)",
					"SJIS.so()(64bit)",
					"T.61.so()(64bit)",
					"TCVN5712-1.so()(64bit)",
					"TIS-620.so()(64bit)",
					"TSCII.so()(64bit)",
					"UHC.so()(64bit)",
					"UNICODE.so()(64bit)",
					"UTF-16.so()(64bit)",
					"UTF-32.so()(64bit)",
					"UTF-7.so()(64bit)",
					"VISCII.so()(64bit)",
					"config(glibc)",
					"ld-linux-x86-64.so.2()(64bit)",
					"ld-linux-x86-64.so.2(GLIBC_2.2.5)(64bit)",
					"ld-linux-x86-64.so.2(GLIBC_2.3)(64bit)",
					"ld-linux-x86-64.so.2(GLIBC_2.4)(64bit)",
					"ldconfig",
					"libBrokenLocale.so.1()(64bit)",
					"libBrokenLocale.so.1(GLIBC_2.2.5)(64bit)",
					"libCNS.so()(64bit)",
					"libGB.so()(64bit)",
					"libISOIR165.so()(64bit)",
					"libJIS.so()(64bit)",
					"libJISX0213.so()(64bit)",
					"libKSC.so()(64bit)",
					"libSegFault.so()(64bit)",
					"libanl.so.1()(64bit)",
					"libanl.so.1(GLIBC_2.2.5)(64bit)",
					"libc.so.6()(64bit)",
					"libc.so.6(GLIBC_2.10)(64bit)",
					"libc.so.6(GLIBC_2.11)(64bit)",
					"libc.so.6(GLIBC_2.12)(64bit)",
					"libc.so.6(GLIBC_2.2.5)(64bit)",
					"libc.so.6(GLIBC_2.2.6)(64bit)",
					"libc.so.6(GLIBC_2.3)(64bit)",
					"libc.so.6(GLIBC_2.3.2)(64bit)",
					"libc.so.6(GLIBC_2.3.3)(64bit)",
					"libc.so.6(GLIBC_2.3.4)(64bit)",
					"libc.so.6(GLIBC_2.4)(64bit)",
					"libc.so.6(GLIBC_2.5)(64bit)",
					"libc.so.6(GLIBC_2.6)(64bit)",
					"libc.so.6(GLIBC_2.7)(64bit)",
					"libc.so.6(GLIBC_2.8)(64bit)",
					"libc.so.6(GLIBC_2.9)(64bit)",
					"libcidn.so.1()(64bit)",
					"libcrypt.so.1()(64bit)",
					"libcrypt.so.1(GLIBC_2.2.5)(64bit)",
					"libdl.so.2()(64bit)",
					"libdl.so.2(GLIBC_2.2.5)(64bit)",
					"libdl.so.2(GLIBC_2.3.3)(64bit)",
					"libdl.so.2(GLIBC_2.3.4)(64bit)",
					"libm.so.6()(64bit)",
					"libm.so.6(GLIBC_2.2.5)(64bit)",
					"libm.so.6(GLIBC_2.4)(64bit)",
					"libmemusage.so()(64bit)",
					"libnsl.so.1()(64bit)",
					"libnsl.so.1(GLIBC_2.2.5)(64bit)",
					"libnss_compat.so.2()(64bit)",
					"libnss_dns.so.2()(64bit)",
					"libnss_files.so.2()(64bit)",
					"libnss_hesiod.so.2()(64bit)",
					"libnss_nis.so.2()(64bit)",
					"libnss_nisplus.so.2()(64bit)",
					"libpcprofile.so()(64bit)",
					"libpthread.so.0()(64bit)",
					"libpthread.so.0(GLIBC_2.11)(64bit)",
					"libpthread.so.0(GLIBC_2.12)(64bit)",
					"libpthread.so.0(GLIBC_2.2.5)(64bit)",
					"libpthread.so.0(GLIBC_2.2.6)(64bit)",
					"libpthread.so.0(GLIBC_2.3.2)(64bit)",
					"libpthread.so.0(GLIBC_2.3.3)(64bit)",
					"libpthread.so.0(GLIBC_2.3.4)(64bit)",
					"libpthread.so.0(GLIBC_2.4)(64bit)",
					"libresolv.so.2()(64bit)",
					"libresolv.so.2(GLIBC_2.2.5)(64bit)",
					"libresolv.so.2(GLIBC_2.3.2)(64bit)",
					"libresolv.so.2(GLIBC_2.9)(64bit)",
					"librt.so.1()(64bit)",
					"librt.so.1(GLIBC_2.2.5)(64bit)",
					"librt.so.1(GLIBC_2.3.3)(64bit)",
					"librt.so.1(GLIBC_2.3.4)(64bit)",
					"librt.so.1(GLIBC_2.4)(64bit)",
					"librt.so.1(GLIBC_2.7)(64bit)",
					"libthread_db.so.1()(64bit)",
					"libthread_db.so.1(GLIBC_2.2.5)(64bit)",
					"libthread_db.so.1(GLIBC_2.3)(64bit)",
					"libthread_db.so.1(GLIBC_2.3.3)(64bit)",
					"libutil.so.1()(64bit)",
					"libutil.so.1(GLIBC_2.2.5)(64bit)",
					"rtld(GNU_HASH)",
					"glibc",
					"glibc(x86-64)",
				},
				Requires: []string{
					"/sbin/ldconfig",
					"/usr/sbin/glibc_post_upgrade.x86_64",
					"basesystem",
					"config(glibc)",
					"glibc-common",
					"ld-linux-x86-64.so.2()(64bit)",
					"ld-linux-x86-64.so.2(GLIBC_2.2.5)(64bit)",
					"ld-linux-x86-64.so.2(GLIBC_2.3)(64bit)",
					"libBrokenLocale.so.1()(64bit)",
					"libCNS.so()(64bit)",
					"libGB.so()(64bit)",
					"libISOIR165.so()(64bit)",
					"libJIS.so()(64bit)",
					"libJISX0213.so()(64bit)",
					"libKSC.so()(64bit)",
					"libanl.so.1()(64bit)",
					"libc.so.6()(64bit)",
					"libc.so.6(GLIBC_2.2.5)(64bit)",
					"libc.so.6(GLIBC_2.3)(64bit)",
					"libc.so.6(GLIBC_2.3.2)(64bit)",
					"libc.so.6(GLIBC_2.3.3)(64bit)",
					"libc.so.6(GLIBC_2.4)(64bit)",
					"libcidn.so.1()(64bit)",
					"libcrypt.so.1()(64bit)",
					"libdl.so.2()(64bit)",
					"libdl.so.2(GLIBC_2.2.5)(64bit)",
					"libfreebl3.so()(64bit)",
					"libfreebl3.so(NSSRAWHASH_3.12.3)(64bit)",
					"libgcc",
					"libm.so.6()(64bit)",
					"libnsl.so.1()(64bit)",
					"libnsl.so.1(GLIBC_2.2.5)(64bit)",
					"libnss_compat.so.2()(64bit)",
					"libnss_dns.so.2()(64bit)",
					"libnss_files.so.2()(64bit)",
					"libnss_hesiod.so.2()(64bit)",
					"libnss_nis.so.2()(64bit)",
					"libnss_nisplus.so.2()(64bit)",
					"libpthread.so.0()(64bit)",
					"libpthread.so.0(GLIBC_2.2.5)(64bit)",
					"libresolv.so.2()(64bit)",
					"libresolv.so.2(GLIBC_2.2.5)(64bit)",
					"libresolv.so.2(GLIBC_2.9)(64bit)",
					"librt.so.1()(64bit)",
					"libthread_db.so.1()(64bit)",
					"libutil.so.1()(64bit)",
					"rpmlib(CompressedFileNames)",
					"rpmlib(FileDigests)",
					"rpmlib(PartialHardlinkSets)",
					"rpmlib(PayloadFilesHavePrefix)",
					"rpmlib(VersionedDependencies)",
					"rpmlib(PayloadIsXz)",
				},
			},
			wantInstalledFiles:     CentOS6GlibcInstalledFiles,
			wantInstalledFileNames: CentOS6GlibcInstalledFileNames,
		},
		{
			name:    "centos8 nodejs",
			pkgName: "nodejs",
			file:    "testdata/centos8-modularitylabel/Packages",
			want: &PackageInfo{
				Epoch:           intRef(1),
				Name:            "nodejs",
				Version:         "10.21.0",
				Release:         "3.module_el8.2.0+391+8da3adc6",
				Arch:            "x86_64",
				Size:            31483781,
				SourceRpm:       "nodejs-10.21.0-3.module_el8.2.0+391+8da3adc6.src.rpm",
				License:         "MIT and ASL 2.0 and ISC and BSD",
				Vendor:          "CentOS",
				Modularitylabel: "nodejs:10:8020020200707141642:6a468ee4",
				Summary:         "JavaScript runtime",
				SigMD5:          "bac7919c2369f944f9da510bbd01370b",
				PGP:             "RSA/SHA256, Tue Jul  7 16:08:24 2020, Key ID 05b555b38483c65d",
				RSAHeader:       "RSA/SHA256, Tue Jul  7 16:08:24 2020, Key ID 05b555b38483c65d",
				DigestAlgorithm: PGPHASHALGO_SHA256,
				InstallTime:     1606911097,
				Provides: []string{
					"bundled(brotli)",
					"bundled(c-ares)",
					"bundled(http-parser)",
					"bundled(icu)",
					"bundled(libuv)",
					"bundled(nghttp2)",
					"bundled(v8)",
					"nodejs",
					"nodejs(abi)",
					"nodejs(abi10)",
					"nodejs(engine)",
					"nodejs(v8-abi)",
					"nodejs(v8-abi6)",
					"nodejs(x86-64)",
					"nodejs-punycode",
					"npm(punycode)",
				},
				Requires: []string{
					"/bin/sh",
					"ca-certificates",
					"libc.so.6()(64bit)",
					"libc.so.6(GLIBC_2.14)(64bit)",
					"libc.so.6(GLIBC_2.15)(64bit)",
					"libc.so.6(GLIBC_2.2.5)(64bit)",
					"libc.so.6(GLIBC_2.28)(64bit)",
					"libc.so.6(GLIBC_2.3)(64bit)",
					"libc.so.6(GLIBC_2.3.2)(64bit)",
					"libc.so.6(GLIBC_2.3.4)(64bit)",
					"libc.so.6(GLIBC_2.4)(64bit)",
					"libc.so.6(GLIBC_2.6)(64bit)",
					"libc.so.6(GLIBC_2.7)(64bit)",
					"libc.so.6(GLIBC_2.9)(64bit)",
					"libcrypto.so.1.1()(64bit)",
					"libcrypto.so.1.1(OPENSSL_1_1_0)(64bit)",
					"libcrypto.so.1.1(OPENSSL_1_1_1)(64bit)",
					"libdl.so.2()(64bit)",
					"libdl.so.2(GLIBC_2.2.5)(64bit)",
					"libgcc_s.so.1()(64bit)",
					"libgcc_s.so.1(GCC_3.0)(64bit)",
					"libgcc_s.so.1(GCC_3.4)(64bit)",
					"libm.so.6()(64bit)",
					"libm.so.6(GLIBC_2.2.5)(64bit)",
					"libpthread.so.0()(64bit)",
					"libpthread.so.0(GLIBC_2.2.5)(64bit)",
					"libpthread.so.0(GLIBC_2.3.2)(64bit)",
					"libpthread.so.0(GLIBC_2.3.3)(64bit)",
					"librt.so.1()(64bit)",
					"librt.so.1(GLIBC_2.2.5)(64bit)",
					"libssl.so.1.1()(64bit)",
					"libssl.so.1.1(OPENSSL_1_1_0)(64bit)",
					"libssl.so.1.1(OPENSSL_1_1_1)(64bit)",
					"libstdc++.so.6()(64bit)",
					"libstdc++.so.6(CXXABI_1.3)(64bit)",
					"libstdc++.so.6(CXXABI_1.3.2)(64bit)",
					"libstdc++.so.6(CXXABI_1.3.5)(64bit)",
					"libstdc++.so.6(CXXABI_1.3.8)(64bit)",
					"libstdc++.so.6(CXXABI_1.3.9)(64bit)",
					"libstdc++.so.6(GLIBCXX_3.4)(64bit)",
					"libstdc++.so.6(GLIBCXX_3.4.11)(64bit)",
					"libstdc++.so.6(GLIBCXX_3.4.14)(64bit)",
					"libstdc++.so.6(GLIBCXX_3.4.15)(64bit)",
					"libstdc++.so.6(GLIBCXX_3.4.18)(64bit)",
					"libstdc++.so.6(GLIBCXX_3.4.20)(64bit)",
					"libstdc++.so.6(GLIBCXX_3.4.21)(64bit)",
					"libstdc++.so.6(GLIBCXX_3.4.9)(64bit)",
					"libz.so.1()(64bit)",
					"npm",
					"rpmlib(CompressedFileNames)",
					"rpmlib(FileDigests)",
					"rpmlib(PayloadFilesHavePrefix)",
					"rpmlib(PayloadIsXz)",
					"rtld(GNU_HASH)",
				},
			},
			wantInstalledFiles:     CentOS8NodejsInstalledFiles,
			wantInstalledFileNames: CentOS8NodejsInstalledFileNames,
		},
		{
			name:    "CBL-Mariner 2.0 curl",
			pkgName: "curl",
			file:    "testdata/cbl-mariner-2.0/rpmdb.sqlite",
			want: &PackageInfo{
				Name:            "curl",
				Version:         "7.76.0",
				Release:         "6.cm2",
				Arch:            "x86_64",
				Size:            326023,
				SourceRpm:       "curl-7.76.0-6.cm2.src.rpm",
				License:         "MIT",
				Vendor:          "Microsoft Corporation",
				Summary:         "An URL retrieval utility and library",
				SigMD5:          "b5f5369ae91df3672fa3338669ec5ca2",
				DigestAlgorithm: PGPHASHALGO_SHA256,
				PGP:             "RSA/SHA256, Thu Jan 27 09:02:11 2022, Key ID 0cd9fed33135ce90",
				RSAHeader:       "RSA/SHA256, Thu Jan 27 09:02:11 2022, Key ID 0cd9fed33135ce90",
				InstallTime:     1643279454,
				Provides: []string{
					"curl",
					"curl(x86-64)",
				},
				Requires: []string{
					"/bin/sh",
					"/sbin/ldconfig",
					"/sbin/ldconfig",
					"curl-libs",
					"krb5",
					"libc.so.6()(64bit)",
					"libc.so.6(GLIBC_2.14)(64bit)",
					"libc.so.6(GLIBC_2.2.5)(64bit)",
					"libc.so.6(GLIBC_2.3)(64bit)",
					"libc.so.6(GLIBC_2.3.4)(64bit)",
					"libc.so.6(GLIBC_2.33)(64bit)",
					"libc.so.6(GLIBC_2.34)(64bit)",
					"libc.so.6(GLIBC_2.4)(64bit)",
					"libc.so.6(GLIBC_2.7)(64bit)",
					"libcurl.so.4()(64bit)",
					"libssh2",
					"libz.so.1()(64bit)",
					"openssl",
					"rpmlib(CompressedFileNames)",
					"rpmlib(FileDigests)",
					"rpmlib(PayloadFilesHavePrefix)",
				},
			},
			wantInstalledFiles:     Mariner2CurlInstalledFiles,
			wantInstalledFileNames: Mariner2CurlInstalledFileNames,
		},
		{
			name:    "Rockylinux 9 bash",
			pkgName: "hostname",
			file:    "testdata/rockylinux-9/rpmdb.sqlite",
			want: &PackageInfo{
				Name:            "hostname",
				Version:         "3.23",
				Release:         "6.el9",
				Arch:            "aarch64",
				Size:            91672,
				SourceRpm:       "hostname-3.23-6.el9.src.rpm",
				License:         "GPLv2+",
				Vendor:          "Rocky Enterprise Software Foundation",
				Summary:         "Utility to set/show the host name or domain name",
				SigMD5:          "8d8cdc55f002f536f30631f92b73d81f",
				DigestAlgorithm: PGPHASHALGO_SHA256,
				PGP:             "", // this is legacy at this point
				RSAHeader:       "RSA/SHA256, Sat May 14 23:43:48 2022, Key ID 702d426d350d275d",
				InstallTime:     1700432743,
				Provides: []string{
					"hostname",
					"hostname(aarch-64)",
				},
				Requires: []string{
					"/bin/sh",
					"/bin/sh",
					"/usr/bin/bash",
					"ld-linux-aarch64.so.1()(64bit)",
					"ld-linux-aarch64.so.1(GLIBC_2.17)(64bit)",
					"libc.so.6()(64bit)",
					"libc.so.6(GLIBC_2.17)(64bit)",
					"libc.so.6(GLIBC_2.33)(64bit)",
					"libc.so.6(GLIBC_2.34)(64bit)",
					"rpmlib(CompressedFileNames)",
					"rpmlib(FileDigests)",
					"rpmlib(PayloadFilesHavePrefix)",
					"rpmlib(PayloadIsZstd)",
					"rtld(GNU_HASH)",
				},
			},
			wantInstalledFiles:     Rockylinux9HostnameFiles,
			wantInstalledFileNames: Rockylinux9HostnameFileNames,
		},
		{
			name:    "libuuid",
			pkgName: "libuuid",
			file:    "testdata/libuuid/Packages",
			want: &PackageInfo{
				Name:            "libuuid",
				Version:         "2.32.1",
				Release:         "42.el8_8",
				Arch:            "x86_64",
				Size:            35104,
				SourceRpm:       "util-linux-2.32.1-42.el8_8.src.rpm",
				License:         "BSD",
				Vendor:          "Red Hat, Inc.",
				Summary:         "Universally unique ID library",
				SigMD5:          "c1e561f13d39aee443a1f00258fba000",
				DigestAlgorithm: PGPHASHALGO_SHA256,
				PGP:             "RSA/SHA256, Mon Apr  3 18:10:39 2023, Key ID 199e2f91fd431d51",
				RSAHeader:       "RSA/SHA256, Mon Apr  3 18:10:39 2023, Key ID 199e2f91fd431d51",
				InstallTime:     1696444673,
				Provides: []string{
					"libuuid",
					"libuuid(x86-64)",
					"libuuid.so.1()(64bit)",
					"libuuid.so.1(UUIDD_PRIVATE)(64bit)",
					"libuuid.so.1(UUID_1.0)(64bit)",
					"libuuid.so.1(UUID_2.20)(64bit)",
					"libuuid.so.1(UUID_2.31)(64bit)",
				},
				Requires: []string{
					"/sbin/ldconfig",
					"/sbin/ldconfig",
					"ld-linux-x86-64.so.2()(64bit)",
					"ld-linux-x86-64.so.2(GLIBC_2.3)(64bit)",
					"libc.so.6()(64bit)",
					"libc.so.6(GLIBC_2.14)(64bit)",
					"libc.so.6(GLIBC_2.2.5)(64bit)",
					"libc.so.6(GLIBC_2.25)(64bit)",
					"libc.so.6(GLIBC_2.28)(64bit)",
					"libc.so.6(GLIBC_2.3)(64bit)",
					"libc.so.6(GLIBC_2.3.4)(64bit)",
					"libc.so.6(GLIBC_2.4)(64bit)",
					"rpmlib(CompressedFileNames)",
					"rpmlib(FileDigests)",
					"rpmlib(PayloadFilesHavePrefix)",
					"rpmlib(PayloadIsXz)",
					"rtld(GNU_HASH)",
				},
			},
			wantInstalledFiles:     LibuuidInstalledFiles,
			wantInstalledFileNames: LibuuidInstalledFileNames,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := Open(tt.file)
			require.NoError(t, err)

			got, err := db.Package(tt.pkgName)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}

			assert.NoError(t, err)

			gotInstalledFiles, err := got.InstalledFiles()
			assert.NoError(t, err)
			assert.Equal(t, tt.wantInstalledFiles, gotInstalledFiles)

			gotInstalledFileNames, err := got.InstalledFileNames()
			assert.NoError(t, err)
			assert.Equal(t, tt.wantInstalledFileNames, gotInstalledFileNames)

			// These fields are tested through InstalledFiles() above
			got.BaseNames = nil
			got.DirIndexes = nil
			got.DirNames = nil
			got.FileSizes = nil
			got.FileDigests = nil
			got.FileModes = nil
			got.FileFlags = nil
			got.UserNames = nil
			got.GroupNames = nil

			assert.Equal(t, tt.want, got)

			err = db.Close()
			require.NoError(t, err)
		})
	}
}

func TestNevra(t *testing.T) {
	blob, err := os.ReadFile("testdata/blob.bin")
	indexEntries, err := headerImport(blob)
	require.NoError(t, err)
	pkg, err := getNEVRA(indexEntries)
	require.NoError(t, err)
	_, err = pkg.InstalledFiles()
	require.Error(t, err)
}
